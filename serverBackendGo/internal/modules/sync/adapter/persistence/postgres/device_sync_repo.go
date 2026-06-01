package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	cfgdomain "github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	profilepostgres "github.com/gis-mdm/server-backend-go/internal/modules/profiles/adapter/persistence/postgres"
	profileapp "github.com/gis-mdm/server-backend-go/internal/modules/profiles/application"
	profiledomain "github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/sync/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/sync/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/sync/port"
	"github.com/gis-mdm/server-backend-go/internal/platform/storage"
	sharedcrypto "github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

type DeviceSyncRepository struct {
	db *sql.DB
}

func NewDeviceSyncRepository(db *sql.DB) *DeviceSyncRepository {
	return &DeviceSyncRepository{db: db}
}

var _ port.SyncRepository = (*DeviceSyncRepository)(nil)

func scanDevice(row *sql.Row) (*domain.DeviceRecord, error) {
	var d domain.DeviceRecord
	var old sql.NullString
	var imei, phone, c1, c2, c3, info sql.NullString
	var enrollmentRoute sql.NullInt64
	err := row.Scan(&d.ID, &d.CustomerID, &d.ConfigurationID, &enrollmentRoute, &d.Number, &old, &imei, &phone,
		&d.LastUpdate, &c1, &c2, &c3, &info)
	if enrollmentRoute.Valid {
		d.EnrollmentRouteID = enrollmentRoute.Int64
	}
	if err != nil {
		return nil, err
	}
	if old.Valid {
		d.OldNumber = &old.String
	}
	if imei.Valid {
		d.IMEI = &imei.String
	}
	if phone.Valid {
		d.Phone = &phone.String
	}
	if c1.Valid {
		d.Custom1 = &c1.String
	}
	if c2.Valid {
		d.Custom2 = &c2.String
	}
	if c3.Valid {
		d.Custom3 = &c3.String
	}
	if info.Valid {
		d.Info = &info.String
	}
	return &d, nil
}

const deviceSelect = `
	SELECT id, customerid, configurationid, enrollment_route_id, number, oldnumber, imei, phone, lastupdate,
	       custom1, custom2, custom3, info
	FROM devices WHERE `

func (r *DeviceSyncRepository) FindByNumber(ctx context.Context, number string) (*domain.DeviceRecord, error) {
	return scanDevice(r.db.QueryRowContext(ctx, deviceSelect+` lower(number) = lower($1)`, number))
}

func (r *DeviceSyncRepository) FindByOldNumber(ctx context.Context, number string) (*domain.DeviceRecord, error) {
	return scanDevice(r.db.QueryRowContext(ctx, deviceSelect+` oldnumber IS NOT NULL AND lower(oldnumber) = lower($1)`, number))
}

func (r *DeviceSyncRepository) CountCustomers(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers`).Scan(&n)
	return n, err
}

func (r *DeviceSyncRepository) CreateOnDemand(ctx context.Context, number string, opts domain.DeviceCreateOptions, defaultCustomerID int64) (*domain.DeviceRecord, error) {
	customerID, err := r.resolveCustomerID(ctx, opts.Customer, defaultCustomerID)
	if err != nil {
		return nil, err
	}
	routeID, legacyConfigID, routeCustomerID, err := r.resolveEnrollmentRoute(ctx, opts.Configuration, customerID)
	if err != nil {
		return nil, err
	}
	if routeCustomerID != customerID {
		return nil, sql.ErrNoRows
	}
	now := time.Now().UnixMilli()
	var treeNodeID sql.NullInt64
	_ = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(er.default_tree_node_id, (
			SELECT n.id FROM device_tree_nodes n
			WHERE n.customerid = $1 AND n.parent_id IS NULL
			ORDER BY n.id LIMIT 1
		))
		FROM enrollment_routes er WHERE er.id = $2`, customerID, routeID).Scan(&treeNodeID)

	var id int64
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO devices (
			number, description, lastupdate, configurationid, customerid, enrolltime,
			tree_node_id, enrollment_route_id, enrollment_state
		)
		VALUES ($1, '', 0, $2, $3, $4, $5, $6, 'enrolled')
		RETURNING id`,
		number, legacyConfigID, customerID, now, nullInt64(treeNodeID), routeID).Scan(&id)
	if err != nil {
		return nil, err
	}
	for _, gname := range opts.Groups {
		gname = strings.TrimSpace(gname)
		if gname == "" {
			continue
		}
		_, _ = r.db.ExecContext(ctx, `
			INSERT INTO devicegroups (deviceid, groupid)
			SELECT $1, g.id FROM groups g
			WHERE g.customerid = $2 AND lower(g.name) = lower($3)
			ON CONFLICT DO NOTHING`, id, customerID, gname)
	}
	return r.FindByNumber(ctx, number)
}

func (r *DeviceSyncRepository) resolveCustomerID(ctx context.Context, customerName string, defaultCustomerID int64) (int64, error) {
	n, err := r.CountCustomers(ctx)
	if err != nil {
		return 0, err
	}
	if n <= 1 {
		return defaultCustomerID, nil
	}
	customerName = strings.TrimSpace(customerName)
	if customerName == "" {
		return 0, sql.ErrNoRows
	}
	var id int64
	err = r.db.QueryRowContext(ctx, `
		SELECT id FROM customers WHERE lower(name) = lower($1)`, customerName).Scan(&id)
	return id, err
}

func (r *DeviceSyncRepository) resolveEnrollmentRoute(ctx context.Context, key string, customerID int64) (routeID, legacyConfigID, routeCustomerID int64, err error) {
	key = strings.TrimSpace(key)
	if key != "" {
		err = r.db.QueryRowContext(ctx, `
			SELECT er.id, COALESCE(er.legacy_configuration_id, er.id), er.customerid
			FROM enrollment_routes er
			WHERE er.qrcodekey IS NOT NULL AND lower(er.qrcodekey) = lower($1)`, key).
			Scan(&routeID, &legacyConfigID, &routeCustomerID)
		if err == nil {
			return routeID, legacyConfigID, routeCustomerID, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, 0, 0, err
		}
		err = r.db.QueryRowContext(ctx, `
			SELECT er.id, COALESCE(er.legacy_configuration_id, er.id), er.customerid
			FROM enrollment_routes er
			WHERE er.customerid = $1 AND lower(er.name) = lower($2)`,
			customerID, key).Scan(&routeID, &legacyConfigID, &routeCustomerID)
		if err == nil {
			return routeID, legacyConfigID, routeCustomerID, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, 0, 0, err
		}
	}
	if key == "" {
		err = r.db.QueryRowContext(ctx, `
			SELECT er.id, COALESCE(er.legacy_configuration_id, er.id), er.customerid
			FROM enrollment_routes er
			WHERE er.customerid = $1
			ORDER BY er.id LIMIT 1`, customerID).
			Scan(&routeID, &legacyConfigID, &routeCustomerID)
		return routeID, legacyConfigID, routeCustomerID, err
	}
	// Legacy fallback before migration backfill
	err = r.db.QueryRowContext(ctx, `
		SELECT id, id, customerid FROM configurations
		WHERE qrcodekey IS NOT NULL AND lower(qrcodekey) = lower($1)`, key).
		Scan(&routeID, &legacyConfigID, &routeCustomerID)
	if err == nil {
		return routeID, legacyConfigID, routeCustomerID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, 0, 0, err
	}
	err = r.db.QueryRowContext(ctx, `
		SELECT id, id, customerid FROM configurations WHERE customerid = $1 AND lower(name) = lower($2)`,
		customerID, key).Scan(&routeID, &legacyConfigID, &routeCustomerID)
	return routeID, legacyConfigID, routeCustomerID, err
}

func (r *DeviceSyncRepository) CompleteMigration(ctx context.Context, deviceID int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE devices SET oldnumber = NULL WHERE id = $1`, deviceID)
	return err
}

func (r *DeviceSyncRepository) TouchLastUpdate(ctx context.Context, deviceID int64) error {
	now := time.Now().UnixMilli()
	_, err := r.db.ExecContext(ctx, `
		UPDATE devices SET lastupdate = $1,
			enrollment_state = CASE WHEN enrollment_state = 'pending' THEN 'enrolled' ELSE 'active' END
		WHERE id = $2`, now, deviceID)
	return err
}

func nullInt64(n sql.NullInt64) any {
	if n.Valid {
		return n.Int64
	}
	return nil
}

func (r *DeviceSyncRepository) UpdateInfo(ctx context.Context, deviceID int64, infoJSON, publicIP string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE devices SET info = $1, infojson = $1::jsonb, publicip = $2, lastupdate = $3 WHERE id = $4`,
		infoJSON, publicIP, time.Now().UnixMilli(), deviceID)
	if err == nil {
		_ = profileapp.RecomputeDeviceRollout(ctx, r.db, deviceID)
	}
	return err
}

func (r *DeviceSyncRepository) UpdateCustomProps(ctx context.Context, deviceID int64, c1, c2, c3 *string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE devices SET custom1 = COALESCE($1, custom1), custom2 = COALESCE($2, custom2),
			custom3 = COALESCE($3, custom3) WHERE id = $4`, c1, c2, c3, deviceID)
	return err
}

func (r *DeviceSyncRepository) SaveApplicationSettings(ctx context.Context, deviceID int64, settings []domain.SyncApplicationSetting) error {
	locks, err := r.policyLocksForDevice(ctx, deviceID)
	if err != nil {
		return err
	}
	for _, s := range settings {
		key := cfgdomain.ApplicationSettingLockKey(s.PackageID, s.Name)
		if locks[key] {
			continue
		}
		_, err := r.db.ExecContext(ctx, `
			DELETE FROM deviceapplicationsettings
			WHERE deviceid = $1 AND applicationpkg = $2 AND name = $3`,
			deviceID, s.PackageID, s.Name)
		if err != nil {
			return err
		}
		_, err = r.db.ExecContext(ctx, `
			INSERT INTO deviceapplicationsettings (deviceid, applicationpkg, name, type, value)
			VALUES ($1, $2, $3, $4, $5)`,
			deviceID, s.PackageID, s.Name, fmt.Sprintf("%d", s.Type), s.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *DeviceSyncRepository) policyLocksForDevice(ctx context.Context, deviceID int64) (map[string]bool, error) {
	var settingsJSON []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT c.settingsjson FROM devices d
		JOIN configurations c ON c.id = d.configurationid
		WHERE d.id = $1`, deviceID).Scan(&settingsJSON)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make(map[string]bool), nil
		}
		return nil, err
	}
	return parsePolicyLocks(settingsJSON), nil
}

func parsePolicyLocks(settingsJSON []byte) map[string]bool {
	out := make(map[string]bool)
	if len(settingsJSON) == 0 {
		return out
	}
	var m map[string]json.RawMessage
	if json.Unmarshal(settingsJSON, &m) != nil {
		return out
	}
	raw, ok := m[cfgdomain.PolicyLocksKey()]
	if !ok {
		return out
	}
	var locks map[string]bool
	if json.Unmarshal(raw, &locks) == nil {
		for k, v := range locks {
			if v {
				out[k] = true
			}
		}
	}
	return out
}

func (r *DeviceSyncRepository) BuildSyncResponse(ctx context.Context, dev domain.DeviceRecord, baseURL, filesDir, cpuArch, mobileName, vendor string) (*domain.SyncResponse, error) {
	if resp, ok, err := r.buildSyncFromArtifact(ctx, dev, baseURL, filesDir, cpuArch, mobileName, vendor); ok {
		if err != nil {
			slog.Warn("sync: artifact path failed, falling through", "deviceId", dev.Number, "err", err)
		} else {
			return resp, nil
		}
	}
	var password, bg, fg, bgi sql.NullString
	var permissive bool
	var settingsJSON []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT password, backgroundcolor, textcolor, backgroundimageurl, permissive, settingsjson
		FROM configurations WHERE id = $1`, dev.ConfigurationID).
		Scan(&password, &bg, &fg, &bgi, &permissive, &settingsJSON)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Info("sync: no legacy configuration, using enrollment route fallback", "deviceId", dev.Number, "configId", dev.ConfigurationID, "routeId", dev.EnrollmentRouteID)
			return r.buildSyncFromEnrollmentRoute(ctx, dev, baseURL, filesDir, cpuArch, mobileName, vendor)
		}
		slog.Error("sync: failed to read configuration", "deviceId", dev.Number, "configId", dev.ConfigurationID, "err", err)
		return nil, err
	}
	resp := &domain.SyncResponse{
		DeviceID:        dev.Number,
		ConfigurationID: dev.ConfigurationID,
		Applications:    []domain.SyncApplication{},
		Files:           []domain.SyncConfigurationFile{},
	}
	if password.Valid && password.String != "" {
		resp.Password = sharedcrypto.MD5UpperHex(password.String)
	}
	if bg.Valid {
		resp.BackgroundColor = &bg.String
	}
	if fg.Valid {
		resp.TextColor = &fg.String
	}
	resp.Permissive = &permissive
	if mobileName != "" {
		resp.AppName = &mobileName
	}
	if vendor != "" {
		resp.Vendor = &vendor
	}
	if dev.OldNumber != nil {
		resp.NewNumber = &dev.Number
	}
	resp.Custom1 = dev.Custom1
	resp.Custom2 = dev.Custom2
	resp.Custom3 = dev.Custom3

	var bgiPtr *string
	if bgi.Valid && bgi.String != "" {
		s := bgi.String
		bgiPtr = &s
	}
	application.ApplyConfigurationPolicy(resp, settingsJSON, bgiPtr)

	rows, err := r.db.QueryContext(ctx, `
		SELECT a.id, a.name, a.pkg, COALESCE(av.version, ''), COALESCE(av.url, ''),
			COALESCE(a.type, 'app'), COALESCE(av.urlarm64, ''), COALESCE(av.urlarmeabi, ''), COALESCE(av.split, false),
			av.versioncode,
			a.runafterinstall, a.runatboot, COALESCE(a.system, false), COALESCE(a.usekiosk, false),
			a.icontext, a.intent,
			COALESCE(ca.showicon, a.showicon, true), COALESCE(ca.remove, false),
			ca.screenorder, ca.keycode, COALESCE(ca.bottom, false), COALESCE(ca.longtap, false),
			COALESCE(cap.skipversioncheck, false),
			uf.filepath
		FROM configurationapplications ca
		JOIN applications a ON a.id = ca.applicationid
		LEFT JOIN applicationversions av ON av.id = ca.applicationversionid
		LEFT JOIN configurationapplicationparameters cap
			ON cap.configurationid = ca.configurationid AND cap.applicationid = ca.applicationid
		LEFT JOIN icons ON icons.id = a.iconid
		LEFT JOIN uploadedfiles uf ON uf.id = icons.fileid
		WHERE ca.configurationid = $1
		ORDER BY ca.screenorder NULLS LAST, a.name`, dev.ConfigurationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	screenIdx := 0
	for rows.Next() {
		var app domain.SyncApplication
		var urlArm64, urlArm, iconPath, iconText, intent sql.NullString
		var split, runAfter, runAtBoot, system, useKiosk, showIcon, remove, bottom, longTap, skipVer bool
		var verCode sql.NullInt64
		var screenOrder, keyCode sql.NullInt64
		if err := rows.Scan(
			&app.ID, &app.Name, &app.Pkg, &app.Version, &app.URL, &app.Type,
			&urlArm64, &urlArm, &split,
			&verCode, &runAfter, &runAtBoot, &system, &useKiosk,
			&iconText, &intent,
			&showIcon, &remove, &screenOrder, &keyCode, &bottom, &longTap, &skipVer,
			&iconPath,
		); err != nil {
			return nil, err
		}
		if split {
			if strings.HasPrefix(cpuArch, "arm64") && urlArm64.Valid && urlArm64.String != "" {
				app.URL = urlArm64.String
			} else if urlArm.Valid && urlArm.String != "" {
				app.URL = urlArm.String
			}
		}
		if verCode.Valid && verCode.Int64 > 0 {
			c := int(verCode.Int64)
			app.Code = &c
		}
		if iconPath.Valid && strings.TrimSpace(iconPath.String) != "" {
			u := storage.BuildPublicURL(baseURL, filesDir, iconPath.String)
			app.Icon = &u
		}
		app.ShowIcon = application.BoolTrueOnly(showIcon)
		if screenOrder.Valid {
			s := int(screenOrder.Int64)
			app.ScreenOrder = &s
		} else {
			screenIdx++
			s := screenIdx
			app.ScreenOrder = &s
		}
		app.UseKiosk = application.BoolTrueOnly(useKiosk)
		app.Remove = application.BoolTrueOnly(remove)
		app.System = application.BoolTrueOnly(system)
		app.RunAfterInstall = application.BoolTrueOnly(runAfter)
		app.RunAtBoot = application.BoolTrueOnly(runAtBoot)
		app.SkipVersion = application.BoolTrueOnly(skipVer)
		app.Bottom = application.BoolTrueOnly(bottom)
		app.LongTap = application.BoolTrueOnly(longTap)
		if iconText.Valid && strings.TrimSpace(iconText.String) != "" {
			t := iconText.String
			app.IconText = &t
		}
		if keyCode.Valid && keyCode.Int64 > 0 {
			k := int(keyCode.Int64)
			app.KeyCode = &k
		}
		if intent.Valid && strings.TrimSpace(intent.String) != "" {
			in := intent.String
			app.Intent = &in
		}
		resp.Applications = append(resp.Applications, app)
	}

	resp.ApplicationSettings = r.mergeApplicationSettings(ctx, dev.ConfigurationID, dev.ID, settingsJSON)

	frows, err := r.db.QueryContext(ctx, `
		SELECT COALESCE(path, ''), COALESCE(externalurl, ''), COALESCE(url, ''), remove
		FROM configurationfiles WHERE configurationid = $1`, dev.ConfigurationID)
	if err != nil {
		return nil, err
	}
	defer frows.Close()
	for frows.Next() {
		var path, extURL, url string
		var remove bool
		if err := frows.Scan(&path, &extURL, &url, &remove); err != nil {
			return nil, err
		}
		f := domain.SyncConfigurationFile{Path: path, Remove: remove}
		if extURL != "" {
			f.URL = extURL
			f.External = true
		} else if path != "" {
			f.URL = storage.BuildPublicURL(baseURL, filesDir, path)
		} else if url != "" {
			f.URL = url
		}
		if f.Path == "" || (!remove && f.URL == "") {
			continue
		}
		resp.Files = append(resp.Files, f)
	}
	return resp, nil
}

func (r *DeviceSyncRepository) buildSyncFromArtifact(ctx context.Context, dev domain.DeviceRecord, baseURL, filesDir, cpuArch, mobileName, vendor string) (*domain.SyncResponse, bool, error) {
	resolved, err := profileapp.ResolveEffectiveProfile(ctx, r.db, dev.ID)
	if err != nil || !resolved.Enabled || resolved.ProfileVersionID <= 0 {
		if dev.EnrollmentRouteID <= 0 {
			return nil, false, nil
		}
		return r.buildSyncFromRouteArtifact(ctx, dev, baseURL, filesDir, cpuArch, mobileName, vendor)
	}
	var raw []byte
	var hash sql.NullString
	err = r.db.QueryRowContext(ctx, `
		SELECT a.artifact_json, a.artifact_hash
		FROM profile_version_artifacts a
		WHERE a.profile_version_id = $1`, resolved.ProfileVersionID).Scan(&raw, &hash)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	routeID := int64(resolved.RouteID)
	if routeID <= 0 && dev.EnrollmentRouteID > 0 {
		routeID = dev.EnrollmentRouteID
	}
	profileID := int64(resolved.ProfileID)
	profileVersionID := int64(resolved.ProfileVersionID)
	_ = profilepostgres.NewAssignmentRepository(r.db).SetDeviceAppliedVersion(ctx, dev.ID, resolved.ProfileVersionID)
	_ = profileapp.RecomputeDeviceRollout(ctx, r.db, dev.ID)
	var artifact profiledomain.ProfileArtifact
	if err := json.Unmarshal(raw, &artifact); err != nil {
		return nil, false, err
	}
	configID := dev.ConfigurationID
	if routeID > 0 {
		configID = routeID
	}
	resp := &domain.SyncResponse{
		DeviceID:        dev.Number,
		ConfigurationID: configID,
		Applications:    []domain.SyncApplication{},
		Files:           []domain.SyncConfigurationFile{},
	}
	if artifact.Password != "" {
		resp.Password = sharedcrypto.MD5UpperHex(artifact.Password)
	}
	if mobileName != "" {
		resp.AppName = &mobileName
	}
	if vendor != "" {
		resp.Vendor = &vendor
	}
	if dev.OldNumber != nil {
		resp.NewNumber = &dev.Number
	}
	resp.Custom1 = dev.Custom1
	resp.Custom2 = dev.Custom2
	resp.Custom3 = dev.Custom3
	pid, pvid := profileID, profileVersionID
	resp.ProfileID = &pid
	resp.ProfileVersionID = &pvid
	if hash.Valid && hash.String != "" {
		rev := hash.String
		resp.ProfileRevision = &rev
	}
	profileapp.ApplyArtifactToSyncResponse(resp, &artifact)
	_ = baseURL
	_ = filesDir
	_ = cpuArch
	resp.ApplicationSettings = r.mergeApplicationSettings(ctx, dev.ConfigurationID, dev.ID, artifact.SettingsJSON)
	return resp, true, nil
}

func (r *DeviceSyncRepository) buildSyncFromRouteArtifact(ctx context.Context, dev domain.DeviceRecord, baseURL, filesDir, cpuArch, mobileName, vendor string) (*domain.SyncResponse, bool, error) {
	var raw []byte
	var routeID, profileID, profileVersionID int64
	var hash sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT er.id, p.id, pv.id, a.artifact_json, a.artifact_hash
		FROM devices d
		JOIN enrollment_routes er ON er.id = d.enrollment_route_id
		JOIN profile_versions pv ON pv.id = er.profile_version_id
		JOIN profiles p ON p.id = pv.profile_id AND COALESCE(p.enabled, true) = true
		JOIN profile_version_artifacts a ON a.profile_version_id = pv.id
		WHERE d.id = $1`, dev.ID).Scan(&routeID, &profileID, &profileVersionID, &raw, &hash)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	_ = profilepostgres.NewAssignmentRepository(r.db).SetDeviceAppliedVersion(ctx, dev.ID, int(profileVersionID))
	_ = profileapp.RecomputeDeviceRollout(ctx, r.db, dev.ID)
	var artifact profiledomain.ProfileArtifact
	if err := json.Unmarshal(raw, &artifact); err != nil {
		return nil, false, err
	}
	configID := dev.ConfigurationID
	if routeID > 0 {
		configID = routeID
	}
	resp := &domain.SyncResponse{
		DeviceID:        dev.Number,
		ConfigurationID: configID,
		Applications:    []domain.SyncApplication{},
		Files:           []domain.SyncConfigurationFile{},
	}
	if artifact.Password != "" {
		resp.Password = sharedcrypto.MD5UpperHex(artifact.Password)
	}
	if mobileName != "" {
		resp.AppName = &mobileName
	}
	if vendor != "" {
		resp.Vendor = &vendor
	}
	if dev.OldNumber != nil {
		resp.NewNumber = &dev.Number
	}
	resp.Custom1 = dev.Custom1
	resp.Custom2 = dev.Custom2
	resp.Custom3 = dev.Custom3
	pid, pvid := profileID, profileVersionID
	resp.ProfileID = &pid
	resp.ProfileVersionID = &pvid
	if hash.Valid && hash.String != "" {
		rev := hash.String
		resp.ProfileRevision = &rev
	}
	profileapp.ApplyArtifactToSyncResponse(resp, &artifact)
	resp.ApplicationSettings = r.mergeApplicationSettings(ctx, dev.ConfigurationID, dev.ID, artifact.SettingsJSON)
	return resp, true, nil
}

// buildSyncFromEnrollmentRoute builds a minimal sync response directly from the enrollment route's
// bootstrap app when no legacy configuration row exists and no profile artifact is available.
// This allows enrollment routes created without a legacy_configuration_id to still complete
// the device enrollment flow successfully.
func (r *DeviceSyncRepository) buildSyncFromEnrollmentRoute(ctx context.Context, dev domain.DeviceRecord, baseURL, filesDir, cpuArch, mobileName, vendor string) (*domain.SyncResponse, error) {
	resp := &domain.SyncResponse{
		DeviceID:        dev.Number,
		ConfigurationID: dev.ConfigurationID,
		Applications:    []domain.SyncApplication{},
		Files:           []domain.SyncConfigurationFile{},
	}
	if mobileName != "" {
		resp.AppName = &mobileName
	}
	if vendor != "" {
		resp.Vendor = &vendor
	}
	if dev.OldNumber != nil {
		resp.NewNumber = &dev.Number
	}
	resp.Custom1 = dev.Custom1
	resp.Custom2 = dev.Custom2
	resp.Custom3 = dev.Custom3
	perm := true
	resp.Permissive = &perm

	// Load the bootstrap app from the enrollment route's mainappid
	if dev.EnrollmentRouteID > 0 {
		var mainAppID sql.NullInt64
		_ = r.db.QueryRowContext(ctx, `
			SELECT mainappid FROM enrollment_routes WHERE id = $1`, dev.EnrollmentRouteID).Scan(&mainAppID)
		if mainAppID.Valid && mainAppID.Int64 > 0 {
			var app domain.SyncApplication
			var urlArm64, urlArm sql.NullString
			var split bool
			var verCode sql.NullInt64
			err := r.db.QueryRowContext(ctx, `
				SELECT a.id, a.name, a.pkg, COALESCE(av.version, ''), COALESCE(av.url, ''),
				       COALESCE(a.type, 'app'), COALESCE(av.urlarm64, ''), COALESCE(av.urlarmeabi, ''),
				       COALESCE(av.split, false), av.versioncode
				FROM applicationversions av
				JOIN applications a ON a.id = av.applicationid
				WHERE av.id = $1`, mainAppID.Int64).
				Scan(&app.ID, &app.Name, &app.Pkg, &app.Version, &app.URL, &app.Type,
					&urlArm64, &urlArm, &split, &verCode)
			if err == nil {
				if split {
					if strings.HasPrefix(cpuArch, "arm64") && urlArm64.Valid && urlArm64.String != "" {
						app.URL = urlArm64.String
					} else if urlArm.Valid && urlArm.String != "" {
						app.URL = urlArm.String
					}
				}
				if verCode.Valid && verCode.Int64 > 0 {
					c := int(verCode.Int64)
					app.Code = &c
				}
				showIcon := true
				app.ShowIcon = &showIcon
				order := 1
				app.ScreenOrder = &order
				resp.Applications = append(resp.Applications, app)
			}
		}
	}
	return resp, nil
}

type appSettingKey struct {
	pkg  string
	name string
}

func (r *DeviceSyncRepository) mergeApplicationSettings(ctx context.Context, configurationID, deviceID int64, settingsJSON []byte) []domain.SyncApplicationSetting {
	locks := parsePolicyLocks(settingsJSON)
	merged := make(map[appSettingKey]domain.SyncApplicationSetting)

	rows, err := r.db.QueryContext(ctx, `
		SELECT COALESCE(a.pkg, ''), s.name, s.type, s.value
		FROM configurationapplicationsettings s
		LEFT JOIN applications a ON a.id = s.applicationid
		WHERE s.configurationid = $1`, configurationID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var s domain.SyncApplicationSetting
			var typ sql.NullString
			if err := rows.Scan(&s.PackageID, &s.Name, &typ, &s.Value); err != nil {
				break
			}
			s.Type = parseSettingType(typ)
			key := appSettingKey{pkg: s.PackageID, name: s.Name}
			lockKey := cfgdomain.ApplicationSettingLockKey(s.PackageID, s.Name)
			if locks[lockKey] {
				s.Readonly = true
			}
			merged[key] = s
		}
	}

	drows, err := r.db.QueryContext(ctx, `
		SELECT applicationpkg, name, type, value
		FROM deviceapplicationsettings WHERE deviceid = $1`, deviceID)
	if err == nil {
		defer drows.Close()
		for drows.Next() {
			var s domain.SyncApplicationSetting
			var typ string
			if err := drows.Scan(&s.PackageID, &s.Name, &typ, &s.Value); err != nil {
				break
			}
			s.Type = parseSettingType(sql.NullString{String: typ, Valid: true})
			key := appSettingKey{pkg: s.PackageID, name: s.Name}
			if existing, ok := merged[key]; ok && existing.Readonly {
				continue
			}
			merged[key] = s
		}
	}

	out := make([]domain.SyncApplicationSetting, 0, len(merged))
	for _, s := range merged {
		out = append(out, s)
	}
	return out
}

func parseSettingType(typ sql.NullString) int {
	if !typ.Valid {
		return 0
	}
	var n int
	if _, err := fmt.Sscanf(typ.String, "%d", &n); err == nil {
		return n
	}
	return 0
}
