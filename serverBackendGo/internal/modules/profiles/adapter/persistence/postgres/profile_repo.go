package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	cfgdomain "github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/port"
)

// ProfileRepository implements port.ProfileRepository.
type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

var _ port.ProfileRepository = (*ProfileRepository)(nil)

var (
	ErrProfileNotFound              = errors.New("profile not found")
	ErrVersionNotFound              = errors.New("profile version not found")
	ErrNotDraftVersion              = errors.New("version is not draft")
	ErrNoPublishedToFork            = errors.New("no published version to fork")
	ErrDuplicateProfileName         = errors.New("duplicate profile name")
	ErrVersionDeleteActivePublished = errors.New("error.profile.version.delete.activePublished")
	ErrVersionDeleteAssigned        = errors.New("error.profile.version.delete.assigned")
	ErrVersionDeleteDevicesTarget   = errors.New("error.profile.version.delete.devicesTarget")
)

func (r *ProfileRepository) List(ctx context.Context, customerID int) ([]domain.ProfileListItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT p.id, p.name, COALESCE(p.description, ''),
		       COALESCE(p.enabled, true),
		       pub.version_number,
		       p.draft_version_id,
		       COALESCE((
		           SELECT COUNT(DISTINCT d.id)::int FROM devices d
		           JOIN enrollment_routes er ON er.id = d.enrollment_route_id
		           JOIN profile_versions pv2 ON pv2.id = er.profile_version_id
		           WHERE pv2.profile_id = p.id
		       ), 0),
		       COALESCE((
		           SELECT COUNT(*)::int FROM enrollment_routes er
		           JOIN profile_versions pv3 ON pv3.id = er.profile_version_id
		           WHERE pv3.profile_id = p.id
		       ), 0)
		FROM profiles p
		LEFT JOIN profile_versions pub ON pub.id = p.published_version_id
		WHERE p.customerid = $1
		ORDER BY lower(p.name)`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ProfileListItem
	for rows.Next() {
		var item domain.ProfileListItem
		var pubVer sql.NullInt64
		var draftID sql.NullInt64
		var routeCount int
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Enabled, &pubVer, &draftID, &item.DeviceCount, &routeCount); err != nil {
			return nil, err
		}
		if pubVer.Valid {
			v := int(pubVer.Int64)
			item.PublishedVersion = &v
		}
		if draftID.Valid && draftID.Int64 > 0 {
			v := int(draftID.Int64)
			item.DraftVersionID = &v
		}
		item.EnrollmentRouteCount = routeCount
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *ProfileRepository) GetMeta(ctx context.Context, customerID, profileID int) (*domain.ProfileMeta, error) {
	var name, desc string
	var enabled bool
	var draftID, pubID sql.NullInt64
	var pubNum sql.NullInt64
	var deviceCount, routeCount int
	err := r.db.QueryRowContext(ctx, `
		SELECT p.name, COALESCE(p.description, ''), COALESCE(p.enabled, true), p.draft_version_id, p.published_version_id,
		       pub.version_number,
		       (SELECT COUNT(DISTINCT d.id)::int FROM devices d
		        JOIN enrollment_routes er ON er.id = d.enrollment_route_id
		        JOIN profile_versions pv2 ON pv2.id = er.profile_version_id
		        WHERE pv2.profile_id = p.id),
		       (SELECT COUNT(*)::int FROM enrollment_routes er
		        JOIN profile_versions pv3 ON pv3.id = er.profile_version_id
		        WHERE pv3.profile_id = p.id)
		FROM profiles p
		LEFT JOIN profile_versions pub ON pub.id = p.published_version_id
		WHERE p.id = $1 AND p.customerid = $2`, profileID, customerID).Scan(
		&name, &desc, &enabled, &draftID, &pubID, &pubNum, &deviceCount, &routeCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	meta := &domain.ProfileMeta{
		ID: profileID, Name: name, Description: desc, Enabled: enabled,
		DeviceCount: deviceCount, EnrollmentRouteCount: routeCount,
	}
	if draftID.Valid && draftID.Int64 > 0 {
		v := int(draftID.Int64)
		meta.DraftVersionID = &v
	}
	if pubID.Valid && pubID.Int64 > 0 {
		v := int(pubID.Int64)
		meta.PublishedVersionID = &v
	}
	if pubNum.Valid {
		v := int(pubNum.Int64)
		meta.PublishedVersion = &v
	}
	return meta, nil
}

func (r *ProfileRepository) GetVersion(ctx context.Context, customerID, profileID, versionID int) (*cfgdomain.Configuration, *domain.VersionMeta, error) {
	ok, err := r.versionBelongsToProfile(ctx, customerID, profileID, versionID)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, nil
	}
	var versionNumber int
	var status string
	var name, desc, password, bg, tc, bgi, qr, base, dfp sql.NullString
	var typ, mainApp, contentApp sql.NullInt64
	var permissive sql.NullBool
	var settings []byte
	err = r.db.QueryRowContext(ctx, `
		SELECT pv.version_number, pv.status,
		       p.name, p.description, pv.type, pv.password, pv.backgroundcolor, pv.textcolor,
		       pv.backgroundimageurl, pv.qrcodekey, pv.baseurl, pv.defaultfilepath,
		       pv.mainappid, pv.contentappid, pv.permissive, pv.settingsjson
		FROM profile_versions pv
		JOIN profiles p ON p.id = pv.profile_id
		WHERE pv.id = $1 AND p.id = $2 AND p.customerid = $3`,
		versionID, profileID, customerID,
	).Scan(
		&versionNumber, &status,
		&name, &desc, &typ, &password, &bg, &tc, &bgi, &qr, &base, &dfp,
		&mainApp, &contentApp, &permissive, &settings,
	)
	if err == sql.ErrNoRows {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	cfg := &cfgdomain.Configuration{ID: &profileID}
	if name.Valid {
		cfg.Name = &name.String
	}
	if desc.Valid {
		cfg.Description = &desc.String
	}
	if typ.Valid {
		t := int(typ.Int64)
		cfg.Type = &t
	}
	if password.Valid {
		cfg.Password = &password.String
	}
	if bg.Valid {
		cfg.BackgroundColor = &bg.String
	}
	if tc.Valid {
		cfg.TextColor = &tc.String
	}
	if bgi.Valid {
		cfg.BackgroundImageURL = &bgi.String
	}
	if qr.Valid {
		cfg.QRCodeKey = &qr.String
	}
	if base.Valid {
		cfg.BaseURL = &base.String
	}
	if dfp.Valid {
		cfg.DefaultFilePath = &dfp.String
	}
	if mainApp.Valid {
		m := int(mainApp.Int64)
		cfg.MainAppID = &m
	}
	if contentApp.Valid {
		m := int(contentApp.Int64)
		cfg.ContentAppID = &m
	}
	if permissive.Valid {
		cfg.Permissive = &permissive.Bool
	}
	if len(settings) > 0 {
		cfg.SetPolicyFromJSON(settings)
	}
	apps, err := r.ListVersionApplications(ctx, customerID, versionID)
	if err != nil {
		return nil, nil, err
	}
	cfg.Applications = apps
	files, err := r.loadFiles(ctx, versionID)
	if err != nil {
		return nil, nil, err
	}
	cfg.Files = files
	settingsRows, err := r.loadAppSettings(ctx, versionID)
	if err != nil {
		return nil, nil, err
	}
	cfg.ApplicationSettings = settingsRows
	meta := &domain.VersionMeta{
		ProfileID: profileID, VersionID: versionID,
		VersionNumber: versionNumber, Status: status,
	}
	return cfg, meta, nil
}

func (r *ProfileRepository) Create(ctx context.Context, customerID int, req domain.CreateRequest) (int, int, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return 0, 0, errors.New("name required")
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	var profileID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO profiles (customerid, name, description)
		VALUES ($1, $2, $3) RETURNING id`,
		customerID, name, nullStr(req.Description),
	).Scan(&profileID)
	if err != nil {
		if strings.Contains(err.Error(), "profiles_name_customer_uidx") {
			return 0, 0, ErrDuplicateProfileName
		}
		return 0, 0, err
	}
	var versionID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO profile_versions (profile_id, version_number, status, settingsjson)
		VALUES ($1, 1, 'draft', '{}'::jsonb) RETURNING id`, profileID,
	).Scan(&versionID)
	if err != nil {
		return 0, 0, err
	}
	if _, err = tx.ExecContext(ctx, `
		UPDATE profiles SET draft_version_id = $1 WHERE id = $2`, versionID, profileID); err != nil {
		return 0, 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, 0, err
	}
	return profileID, versionID, nil
}

func (r *ProfileRepository) EnsureDraft(ctx context.Context, customerID, profileID int) (int, error) {
	meta, err := r.GetMeta(ctx, customerID, profileID)
	if err != nil || meta == nil {
		return 0, ErrProfileNotFound
	}
	if meta.DraftVersionID != nil && *meta.DraftVersionID > 0 {
		var status string
		err := r.db.QueryRowContext(ctx, `
			SELECT status FROM profile_versions WHERE id = $1`, *meta.DraftVersionID).Scan(&status)
		if err == nil && status == "draft" {
			return *meta.DraftVersionID, nil
		}
	}
	if meta.PublishedVersionID == nil || *meta.PublishedVersionID == 0 {
		return 0, ErrNoPublishedToFork
	}
	return r.forkDraftFromVersion(ctx, profileID, *meta.PublishedVersionID)
}

func (r *ProfileRepository) forkDraftFromVersion(ctx context.Context, profileID, sourceVersionID int) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var nextNum int
	if err := tx.QueryRowContext(ctx, `
		SELECT COALESCE(MAX(version_number), 0) + 1 FROM profile_versions WHERE profile_id = $1`,
		profileID,
	).Scan(&nextNum); err != nil {
		return 0, err
	}
	var newID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO profile_versions (
			profile_id, version_number, status, type, password,
			backgroundcolor, textcolor, backgroundimageurl, qrcodekey, baseurl,
			defaultfilepath, mainappid, contentappid, permissive, settingsjson
		)
		SELECT profile_id, $2, 'draft', type, password,
		       backgroundcolor, textcolor, backgroundimageurl, qrcodekey, baseurl,
		       defaultfilepath, mainappid, contentappid, permissive, settingsjson
		FROM profile_versions WHERE id = $1
		RETURNING id`, sourceVersionID, nextNum,
	).Scan(&newID)
	if err != nil {
		return 0, err
	}
	for _, stmt := range []string{
		`INSERT INTO profile_version_applications (
			profile_version_id, applicationid, applicationversionid, action, showicon,
			screenorder, keycode, bottom, remove, longtap
		)
		SELECT $1, applicationid, applicationversionid, action, showicon,
		       screenorder, keycode, bottom, remove, longtap
		FROM profile_version_applications WHERE profile_version_id = $2`,
		`INSERT INTO profile_version_files (profile_version_id, path, externalurl, url, remove)
		SELECT $1, path, externalurl, url, remove
		FROM profile_version_files WHERE profile_version_id = $2`,
		`INSERT INTO profile_version_application_settings (profile_version_id, applicationid, name, type, value)
		SELECT $1, applicationid, name, type, value
		FROM profile_version_application_settings WHERE profile_version_id = $2`,
		`INSERT INTO profile_version_application_parameters (profile_version_id, applicationid, name, value)
		SELECT $1, applicationid, name, value
		FROM profile_version_application_parameters WHERE profile_version_id = $2`,
	} {
		if _, err = tx.ExecContext(ctx, stmt, newID, sourceVersionID); err != nil {
			return 0, err
		}
	}
	if _, err = tx.ExecContext(ctx, `
		UPDATE profiles SET draft_version_id = $1 WHERE id = $2`, newID, profileID); err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return newID, nil
}

func (r *ProfileRepository) SaveDraft(ctx context.Context, customerID, profileID, versionID int, cfg cfgdomain.Configuration) error {
	ok, err := r.versionBelongsToProfile(ctx, customerID, profileID, versionID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrVersionNotFound
	}
	var status string
	if err := r.db.QueryRowContext(ctx, `
		SELECT status FROM profile_versions WHERE id = $1`, versionID).Scan(&status); err != nil {
		return err
	}
	if status != "draft" {
		return ErrNotDraftVersion
	}
	settings, err := cfg.BuildSettingsJSON()
	if err != nil {
		return err
	}
	typ := 0
	if cfg.Type != nil {
		typ = *cfg.Type
	}
	permissive := false
	if cfg.Permissive != nil {
		permissive = *cfg.Permissive
	}
	name := ""
	if cfg.Name != nil {
		name = strings.TrimSpace(*cfg.Name)
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, `
		UPDATE profiles SET name = $1, description = $2 WHERE id = $3 AND customerid = $4`,
		name, nullStr(cfg.Description), profileID, customerID,
	); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `
		UPDATE profile_versions SET
			type=$1, password=$2, backgroundcolor=$3, textcolor=$4,
			backgroundimageurl=$5, qrcodekey=$6, baseurl=$7, defaultfilepath=$8,
			mainappid=$9, contentappid=$10, permissive=$11, settingsjson=$12
		WHERE id=$13`,
		typ, nullStr(cfg.Password), nullStr(cfg.BackgroundColor), nullStr(cfg.TextColor),
		nullStr(cfg.BackgroundImageURL), nullStr(cfg.QRCodeKey), nullStr(cfg.BaseURL),
		nullStr(cfg.DefaultFilePath), nullInt(cfg.MainAppID), nullInt(cfg.ContentAppID),
		permissive, settings, versionID,
	); err != nil {
		return err
	}
	for _, table := range []string{
		"profile_version_application_parameters",
		"profile_version_applications",
		"profile_version_files",
		"profile_version_application_settings",
	} {
		if _, err = tx.ExecContext(ctx, fmt.Sprintf(`DELETE FROM %s WHERE profile_version_id=$1`, table), versionID); err != nil {
			return err
		}
	}
	for _, app := range cfg.Applications {
		appID := app.ID
		if app.ApplicationID > 0 {
			appID = app.ApplicationID
		}
		if appID <= 0 {
			continue
		}
		verID := 0
		if app.UsedVersionID != nil {
			verID = *app.UsedVersionID
		}
		if verID <= 0 {
			var latest int
			if err := tx.QueryRowContext(ctx, `
				SELECT av.id FROM applicationversions av
				JOIN applications a ON a.id = av.applicationid
				WHERE av.applicationid = $1 AND (a.customerid = $2 OR a.common = TRUE)
				ORDER BY av.versioncode DESC, av.id DESC LIMIT 1`,
				appID, customerID).Scan(&latest); err != nil || latest <= 0 {
				continue
			}
			verID = latest
		}
		action := 1
		if app.Action != nil {
			action = *app.Action
		}
		showIcon := true
		if app.ShowIcon != nil {
			showIcon = *app.ShowIcon
		}
		bottom, remove, longTap := false, false, false
		if app.Bottom != nil {
			bottom = *app.Bottom
		}
		if app.Remove != nil {
			remove = *app.Remove
		}
		if app.LongTap != nil {
			longTap = *app.LongTap
		}
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO profile_version_applications (
				profile_version_id, applicationid, applicationversionid,
				action, showicon, screenorder, keycode, bottom, remove, longtap
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
			versionID, appID, verID, action, showIcon,
			nullInt(app.ScreenOrder), nullInt(app.KeyCode), bottom, remove, longTap,
		); err != nil {
			return err
		}
		if app.SkipVersionCheck != nil {
			val := "false"
			if *app.SkipVersionCheck {
				val = "true"
			}
			if _, err = tx.ExecContext(ctx, `
				INSERT INTO profile_version_application_parameters (profile_version_id, applicationid, name, value)
				VALUES ($1,$2,'skipVersionCheck',$3)
				ON CONFLICT (profile_version_id, applicationid, name) DO UPDATE SET value = $3`,
				versionID, appID, val,
			); err != nil {
				return err
			}
		}
	}
	for _, f := range cfg.Files {
		if f.Remove != nil && *f.Remove {
			continue
		}
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO profile_version_files (profile_version_id, path, externalurl, url, remove)
			VALUES ($1,$2,$3,$4,$5)`,
			versionID, nullStr(f.Path), nullStr(f.ExternalURL), nullStr(f.URL), false,
		); err != nil {
			return err
		}
	}
	for _, s := range cfg.ApplicationSettings {
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO profile_version_application_settings (profile_version_id, applicationid, name, type, value)
			VALUES ($1,$2,$3,$4,$5)`,
			versionID, nullInt(s.ApplicationID), nullStr(s.Name), nullStr(s.Type), nullStr(s.Value),
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *ProfileRepository) ListVersionApplications(ctx context.Context, customerID, versionID int) ([]cfgdomain.ConfigurationApplication, error) {
	_ = customerID
	rows, err := r.db.QueryContext(ctx, `
		SELECT a.id, a.name, a.pkg, a.type,
		       pva.applicationversionid, pva.action, pva.showicon,
		       pva.screenorder, pva.keycode, pva.bottom,
		       pva.remove, pva.longtap,
		       COALESCE(pap.value = 'true', false),
		       av.version, av.versioncode, av.url,
		       (SELECT av2.id FROM applicationversions av2
		        WHERE av2.applicationid = a.id ORDER BY av2.versioncode DESC, av2.id DESC LIMIT 1) AS latest_version
		FROM profile_version_applications pva
		JOIN applications a ON a.id = pva.applicationid
		LEFT JOIN applicationversions av ON av.id = pva.applicationversionid
		LEFT JOIN profile_version_application_parameters pap
		       ON pap.profile_version_id = pva.profile_version_id
		      AND pap.applicationid = pva.applicationid
		      AND pap.name = 'skipVersionCheck'
		WHERE pva.profile_version_id = $1
		ORDER BY lower(a.name)`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]cfgdomain.ConfigurationApplication, 0)
	for rows.Next() {
		var app cfgdomain.ConfigurationApplication
		var name, pkg, typ, ver, url sql.NullString
		var usedVer, latest, action, verCode sql.NullInt64
		var showIcon, bottom, remove, longTap, skipVer sql.NullBool
		var screenOrder, keyCode sql.NullInt64
		if err := rows.Scan(
			&app.ID, &name, &pkg, &typ, &usedVer, &action, &showIcon,
			&screenOrder, &keyCode, &bottom, &remove, &longTap, &skipVer,
			&ver, &verCode, &url, &latest,
		); err != nil {
			return nil, err
		}
		if name.Valid {
			app.Name = &name.String
		}
		if pkg.Valid {
			app.Pkg = &pkg.String
		}
		if typ.Valid {
			app.Type = &typ.String
		}
		if usedVer.Valid && usedVer.Int64 > 0 {
			v := int(usedVer.Int64)
			app.UsedVersionID = &v
		}
		if action.Valid {
			a := int(action.Int64)
			app.Action = &a
		}
		if showIcon.Valid {
			app.ShowIcon = &showIcon.Bool
		}
		if screenOrder.Valid {
			s := int(screenOrder.Int64)
			app.ScreenOrder = &s
		}
		if keyCode.Valid {
			k := int(keyCode.Int64)
			app.KeyCode = &k
		}
		if bottom.Valid {
			app.Bottom = &bottom.Bool
		}
		if remove.Valid {
			app.Remove = &remove.Bool
		}
		if longTap.Valid {
			app.LongTap = &longTap.Bool
		}
		if skipVer.Valid {
			app.SkipVersionCheck = &skipVer.Bool
		}
		if ver.Valid {
			app.Version = &ver.String
		}
		if verCode.Valid {
			v := int(verCode.Int64)
			app.VersionCode = &v
		}
		if url.Valid {
			app.URL = &url.String
		}
		if latest.Valid {
			l := int(latest.Int64)
			app.LatestVersion = &l
		}
		app.ApplicationID = app.ID
		out = append(out, app)
	}
	return out, rows.Err()
}

func (r *ProfileRepository) versionBelongsToProfile(ctx context.Context, customerID, profileID, versionID int) (bool, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `
		SELECT 1 FROM profile_versions pv
		JOIN profiles p ON p.id = pv.profile_id
		WHERE pv.id = $1 AND p.id = $2 AND p.customerid = $3`,
		versionID, profileID, customerID,
	).Scan(&n)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *ProfileRepository) loadFiles(ctx context.Context, versionID int) ([]cfgdomain.ConfigurationFile, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, path, externalurl, url, remove
		FROM profile_version_files WHERE profile_version_id = $1`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []cfgdomain.ConfigurationFile
	for rows.Next() {
		var f cfgdomain.ConfigurationFile
		var id int
		var path, ext, url sql.NullString
		var remove bool
		if err := rows.Scan(&id, &path, &ext, &url, &remove); err != nil {
			return nil, err
		}
		f.ID = &id
		if path.Valid {
			f.Path = &path.String
		}
		if ext.Valid {
			f.ExternalURL = &ext.String
		}
		if url.Valid {
			f.URL = &url.String
		}
		f.Remove = &remove
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *ProfileRepository) loadAppSettings(ctx context.Context, versionID int) ([]cfgdomain.ConfigurationApplicationSetting, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.applicationid, a.name, s.name, s.type, s.value
		FROM profile_version_application_settings s
		LEFT JOIN applications a ON a.id = s.applicationid
		WHERE s.profile_version_id = $1`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []cfgdomain.ConfigurationApplicationSetting
	for rows.Next() {
		var s cfgdomain.ConfigurationApplicationSetting
		var id int
		var appID sql.NullInt64
		var appName, name, typ, val sql.NullString
		if err := rows.Scan(&id, &appID, &appName, &name, &typ, &val); err != nil {
			return nil, err
		}
		s.ID = &id
		if appID.Valid {
			a := int(appID.Int64)
			s.ApplicationID = &a
		}
		if appName.Valid {
			s.ApplicationName = &appName.String
		}
		if name.Valid {
			s.Name = &name.String
		}
		if typ.Valid {
			s.Type = &typ.String
		}
		if val.Valid {
			s.Value = &val.String
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func nullStr(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}

func publishedByOrNull(userID int) any {
	if userID <= 0 {
		return nil
	}
	return userID
}

func nullInt(n *int) any {
	if n == nil || *n == 0 {
		return nil
	}
	return *n
}

func (r *ProfileRepository) CountImpact(ctx context.Context, customerID, profileID int) (int, int, error) {
	var devices, routes int
	err := r.db.QueryRowContext(ctx, `
		SELECT
			(SELECT COUNT(DISTINCT d.id)::int FROM devices d
			 JOIN enrollment_routes er ON er.id = d.enrollment_route_id
			 JOIN profile_versions pv ON pv.id = er.profile_version_id
			 WHERE pv.profile_id = p.id),
			(SELECT COUNT(*)::int FROM enrollment_routes er
			 JOIN profile_versions pv ON pv.id = er.profile_version_id
			 WHERE pv.profile_id = p.id)
		FROM profiles p WHERE p.id = $1 AND p.customerid = $2`, profileID, customerID).
		Scan(&devices, &routes)
	return devices, routes, err
}

func (r *ProfileRepository) PublishVersion(ctx context.Context, customerID, profileID, versionID int, artifactJSON []byte, artifactHash string, publishedBy int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := r.PublishVersionTx(ctx, tx, customerID, profileID, versionID, artifactJSON, artifactHash, publishedBy); err != nil {
		return err
	}
	return tx.Commit()
}

// PublishVersionTx publishes a draft inside an existing transaction (020 assignment bump).
func (r *ProfileRepository) PublishVersionTx(ctx context.Context, tx *sql.Tx, customerID, profileID, versionID int, artifactJSON []byte, artifactHash string, publishedBy int) error {
	ok, err := r.versionBelongsToProfileTx(ctx, tx, customerID, profileID, versionID)
	if err != nil || !ok {
		return ErrVersionNotFound
	}
	var status string
	if err := tx.QueryRowContext(ctx, `SELECT status FROM profile_versions WHERE id = $1`, versionID).Scan(&status); err != nil {
		return err
	}
	if status != "draft" {
		return ErrNotDraftVersion
	}
	if _, err = tx.ExecContext(ctx, `
		UPDATE profile_versions SET status = 'archived'
		WHERE profile_id = $1 AND status = 'published'`, profileID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `
		UPDATE profile_versions SET status = 'published', published_at = NOW(), published_by = $1
		WHERE id = $2`, publishedByOrNull(publishedBy), versionID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `
		UPDATE profiles SET published_version_id = $1, draft_version_id = NULL
		WHERE id = $2 AND customerid = $3`, versionID, profileID, customerID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `
		INSERT INTO profile_version_artifacts (profile_version_id, artifact_json, artifact_hash, compiled_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (profile_version_id) DO UPDATE SET
			artifact_json = EXCLUDED.artifact_json,
			artifact_hash = EXCLUDED.artifact_hash,
			compiled_at = NOW()`, versionID, artifactJSON, artifactHash); err != nil {
		return err
	}
	return nil
}

func (r *ProfileRepository) versionBelongsToProfileTx(ctx context.Context, tx *sql.Tx, customerID, profileID, versionID int) (bool, error) {
	var ok bool
	err := tx.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM profile_versions pv
			JOIN profiles p ON p.id = pv.profile_id
			WHERE pv.id = $1 AND p.id = $2 AND p.customerid = $3
		)`, versionID, profileID, customerID).Scan(&ok)
	return ok, err
}

func (r *ProfileRepository) VersionDeleteEligibility(ctx context.Context, customerID, profileID, versionID int) (activePublished, assigned, devicesTarget bool, err error) {
	ok, err := r.versionBelongsToProfile(ctx, customerID, profileID, versionID)
	if err != nil || !ok {
		return false, false, false, ErrVersionNotFound
	}
	var pubID sql.NullInt64
	if err = r.db.QueryRowContext(ctx, `
		SELECT published_version_id FROM profiles WHERE id = $1 AND customerid = $2`, profileID, customerID).Scan(&pubID); err != nil {
		if err == sql.ErrNoRows {
			return false, false, false, ErrProfileNotFound
		}
		return false, false, false, err
	}
	if pubID.Valid && int(pubID.Int64) == versionID {
		activePublished = true
	}
	if err = r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM profile_tree_assignments
			WHERE customerid = $1 AND profile_id = $2 AND profile_version_id = $3
		)`, customerID, profileID, versionID).Scan(&assigned); err != nil {
		return false, false, false, err
	}
	if err = r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM devices WHERE customerid = $1 AND target_profile_version_id = $2
		)`, customerID, versionID).Scan(&devicesTarget); err != nil {
		return false, false, false, err
	}
	return activePublished, assigned, devicesTarget, nil
}

func (r *ProfileRepository) DeleteVersion(ctx context.Context, customerID, profileID, versionID int) error {
	activePublished, assigned, devicesTarget, err := r.VersionDeleteEligibility(ctx, customerID, profileID, versionID)
	if err != nil {
		return err
	}
	if activePublished {
		return ErrVersionDeleteActivePublished
	}
	if assigned {
		return ErrVersionDeleteAssigned
	}
	if devicesTarget {
		return ErrVersionDeleteDevicesTarget
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `
		UPDATE profiles SET draft_version_id = NULL
		WHERE id = $1 AND customerid = $2 AND draft_version_id = $3`, profileID, customerID, versionID); err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, `
		DELETE FROM profile_versions pv
		USING profiles p
		WHERE pv.id = $1 AND pv.profile_id = $2 AND p.id = pv.profile_id AND p.customerid = $3`,
		versionID, profileID, customerID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrVersionNotFound
	}
	return tx.Commit()
}

func (r *ProfileRepository) InsertDomainEvent(ctx context.Context, eventType, aggregateID string, payload []byte) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO domain_events (event_type, aggregate_id, payload)
		VALUES ($1, $2, $3)`, eventType, aggregateID, payload)
	return err
}

// ListVersions returns all versions for a profile (018 version navigation).
func (r *ProfileRepository) ListVersions(ctx context.Context, customerID, profileID int) ([]domain.VersionListItem, error) {
	var exists int
	if err := r.db.QueryRowContext(ctx, `
		SELECT 1 FROM profiles WHERE id = $1 AND customerid = $2`, profileID, customerID).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT pv.id, pv.version_number, pv.status,
		       pv.published_at, pv.created_at
		FROM profile_versions pv
		WHERE pv.profile_id = $1
		ORDER BY pv.version_number DESC`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.VersionListItem
	for rows.Next() {
		var item domain.VersionListItem
		var pubAt sql.NullTime
		var created time.Time
		if err := rows.Scan(&item.VersionID, &item.VersionNumber, &item.Status, &pubAt, &created); err != nil {
			return nil, err
		}
		item.CreatedAt = created.UTC().Format(time.RFC3339)
		if pubAt.Valid {
			s := pubAt.Time.UTC().Format(time.RFC3339)
			item.PublishedAt = &s
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

// ForkDraftFromPublished creates a draft copy from a published version id.
func (r *ProfileRepository) ForkDraftFromPublished(ctx context.Context, customerID, profileID, sourceVersionID int) (int, error) {
	ok, err := r.versionBelongsToProfile(ctx, customerID, profileID, sourceVersionID)
	if err != nil || !ok {
		return 0, ErrVersionNotFound
	}
	var status string
	if err := r.db.QueryRowContext(ctx, `
		SELECT status FROM profile_versions WHERE id = $1`, sourceVersionID).Scan(&status); err != nil {
		return 0, err
	}
	if status != "published" {
		return 0, ErrNotDraftVersion
	}
	return r.forkDraftFromVersion(ctx, profileID, sourceVersionID)
}
