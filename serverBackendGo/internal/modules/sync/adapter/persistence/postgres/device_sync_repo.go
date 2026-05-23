package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	cfgdomain "github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
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
	err := row.Scan(&d.ID, &d.CustomerID, &d.ConfigurationID, &d.Number, &old, &imei, &phone,
		&d.LastUpdate, &c1, &c2, &c3, &info)
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
	SELECT id, customerid, configurationid, number, oldnumber, imei, phone, lastupdate,
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
	configID, configCustomerID, err := r.resolveConfiguration(ctx, opts.Configuration, customerID)
	if err != nil {
		return nil, err
	}
	if configCustomerID != customerID {
		return nil, sql.ErrNoRows
	}
	now := time.Now().UnixMilli()
	var id int64
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO devices (number, description, lastupdate, configurationid, customerid, enrolltime)
		VALUES ($1, '', 0, $2, $3, $4)
		RETURNING id`, number, configID, customerID, now).Scan(&id)
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

func (r *DeviceSyncRepository) resolveConfiguration(ctx context.Context, key string, customerID int64) (configID, configCustomerID int64, err error) {
	key = strings.TrimSpace(key)
	if key == "" {
		err = r.db.QueryRowContext(ctx, `
			SELECT id, customerid FROM configurations WHERE customerid = $1 ORDER BY id LIMIT 1`, customerID).
			Scan(&configID, &configCustomerID)
		return configID, configCustomerID, err
	}
	err = r.db.QueryRowContext(ctx, `
		SELECT id, customerid FROM configurations
		WHERE qrcodekey IS NOT NULL AND lower(qrcodekey) = lower($1)`, key).
		Scan(&configID, &configCustomerID)
	if err == nil {
		return configID, configCustomerID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, 0, err
	}
	err = r.db.QueryRowContext(ctx, `
		SELECT id, customerid FROM configurations WHERE customerid = $1 AND lower(name) = lower($2)`,
		customerID, key).Scan(&configID, &configCustomerID)
	return configID, configCustomerID, err
}

func (r *DeviceSyncRepository) CompleteMigration(ctx context.Context, deviceID int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE devices SET oldnumber = NULL WHERE id = $1`, deviceID)
	return err
}

func (r *DeviceSyncRepository) TouchLastUpdate(ctx context.Context, deviceID int64) error {
	now := time.Now().UnixMilli()
	_, err := r.db.ExecContext(ctx, `UPDATE devices SET lastupdate = $1 WHERE id = $2`, now, deviceID)
	return err
}

func (r *DeviceSyncRepository) UpdateInfo(ctx context.Context, deviceID int64, infoJSON, publicIP string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE devices SET info = $1, publicip = $2, lastupdate = $3 WHERE id = $4`,
		infoJSON, publicIP, time.Now().UnixMilli(), deviceID)
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
	var password, bg, fg, bgi sql.NullString
	var permissive bool
	var settingsJSON []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT password, backgroundcolor, textcolor, backgroundimageurl, permissive, settingsjson
		FROM configurations WHERE id = $1`, dev.ConfigurationID).
		Scan(&password, &bg, &fg, &bgi, &permissive, &settingsJSON)
	if err != nil {
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
			COALESCE(a.type, 'app'), COALESCE(av.urlarm64, ''), COALESCE(av.urlarmeabi, ''), COALESCE(av.split, false)
		FROM configurationapplications ca
		JOIN applications a ON a.id = ca.applicationid
		LEFT JOIN applicationversions av ON av.id = ca.applicationversionid
		WHERE ca.configurationid = $1
		ORDER BY ca.screenorder NULLS LAST, a.name`, dev.ConfigurationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var app domain.SyncApplication
		var urlArm64, urlArm sql.NullString
		var split bool
		if err := rows.Scan(&app.ID, &app.Name, &app.Pkg, &app.Version, &app.URL, &app.Type, &urlArm64, &urlArm, &split); err != nil {
			return nil, err
		}
		if split {
			if strings.HasPrefix(cpuArch, "arm64") && urlArm64.Valid && urlArm64.String != "" {
				app.URL = urlArm64.String
			} else if urlArm.Valid && urlArm.String != "" {
				app.URL = urlArm.String
			}
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
