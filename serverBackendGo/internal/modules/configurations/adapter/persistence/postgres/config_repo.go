package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/configurations/port"
)

// ConfigRepository implements port.ConfigRepository.
type ConfigRepository struct {
	db *sql.DB
}

func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

var _ port.ConfigRepository = (*ConfigRepository)(nil)

func (r *ConfigRepository) ListByCustomer(ctx context.Context, customerID int) ([]domain.LookupItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name FROM configurations WHERE customerid = $1 ORDER BY lower(name)`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.LookupItem
	for rows.Next() {
		var item domain.LookupItem
		var name string
		if err := rows.Scan(&item.ID, &name); err != nil {
			return nil, err
		}
		item.Name = &name
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *ConfigRepository) Search(ctx context.Context, customerID int) ([]domain.Configuration, error) {
	return r.searchQuery(ctx, customerID, "")
}

func (r *ConfigRepository) SearchByValue(ctx context.Context, customerID int, value string) ([]domain.Configuration, error) {
	return r.searchQuery(ctx, customerID, strings.TrimSpace(value))
}

func (r *ConfigRepository) searchQuery(ctx context.Context, customerID int, filter string) ([]domain.Configuration, error) {
	q := `
		SELECT c.id, c.name, c.description, c.type,
		       (SELECT COUNT(*)::int FROM devices d WHERE d.configurationid = c.id) AS device_count,
		       c.qrcodekey, c.mainappid
		FROM configurations c
		WHERE c.customerid = $1`
	args := []any{customerID}
	if filter != "" {
		q += ` AND (c.name ILIKE $2 OR COALESCE(c.description,'') ILIKE $2)`
		args = append(args, "%"+filter+"%")
	}
	q += ` ORDER BY lower(c.name)`
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Configuration
	for rows.Next() {
		var c domain.Configuration
		var id, typ, dc int
		var name, desc, qr sql.NullString
		var mainApp sql.NullInt64
		if err := rows.Scan(&id, &name, &desc, &typ, &dc, &qr, &mainApp); err != nil {
			return nil, err
		}
		c.ID = &id
		if name.Valid {
			c.Name = &name.String
		}
		if desc.Valid {
			c.Description = &desc.String
		}
		c.Type = &typ
		c.DeviceCount = &dc
		if qr.Valid && strings.TrimSpace(qr.String) != "" {
			c.QRCodeKey = &qr.String
		}
		if mainApp.Valid && mainApp.Int64 > 0 {
			m := int(mainApp.Int64)
			c.MainAppID = &m
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *ConfigRepository) GetByName(ctx context.Context, customerID int, name string) (*domain.Configuration, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM configurations
		WHERE customerid = $1 AND lower(name) = lower($2)`, customerID, name).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, customerID, id)
}

func (r *ConfigRepository) GetByID(ctx context.Context, customerID, id int) (*domain.Configuration, error) {
	var name, desc, password, bg, tc, bgi, qr, base, dfp sql.NullString
	var typ, mainApp, contentApp sql.NullInt64
	var permissive sql.NullBool
	var settings []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT name, description, type, password, backgroundcolor, textcolor,
		       backgroundimageurl, qrcodekey, baseurl, defaultfilepath,
		       mainappid, contentappid, permissive, settingsjson
		FROM configurations
		WHERE id = $1 AND customerid = $2`, id, customerID).Scan(
		&name, &desc, &typ, &password, &bg, &tc, &bgi, &qr, &base, &dfp,
		&mainApp, &contentApp, &permissive, &settings,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	cfg := &domain.Configuration{ID: &id}
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
	cfg.ID = &id
	apps, err := r.ListConfigurationApplications(ctx, customerID, id)
	if err != nil {
		return nil, err
	}
	cfg.Applications = apps
	files, err := r.loadFiles(ctx, id)
	if err != nil {
		return nil, err
	}
	cfg.Files = files
	settingsRows, err := r.loadAppSettings(ctx, id)
	if err != nil {
		return nil, err
	}
	cfg.ApplicationSettings = settingsRows
	return cfg, nil
}

func (r *ConfigRepository) CountDevicesUsing(ctx context.Context, configurationID int) (int64, error) {
	var n int64
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM devices WHERE configurationid = $1`, configurationID).Scan(&n)
	return n, err
}

func (r *ConfigRepository) Insert(ctx context.Context, customerID int, cfg domain.Configuration) (int, error) {
	return r.save(ctx, customerID, 0, cfg)
}

func (r *ConfigRepository) Update(ctx context.Context, customerID int, cfg domain.Configuration) error {
	if cfg.ID == nil || *cfg.ID == 0 {
		return errors.New("configuration id required")
	}
	_, err := r.save(ctx, customerID, *cfg.ID, cfg)
	return err
}

func (r *ConfigRepository) save(ctx context.Context, customerID, id int, cfg domain.Configuration) (int, error) {
	settings, err := cfg.BuildSettingsJSON()
	if err != nil {
		return 0, err
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
		return 0, err
	}
	defer tx.Rollback()

	if id == 0 {
		err = tx.QueryRowContext(ctx, `
			INSERT INTO configurations (
				name, description, customerid, type, password,
				backgroundcolor, textcolor, backgroundimageurl, qrcodekey, baseurl,
				defaultfilepath, mainappid, contentappid, permissive, settingsjson
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
			RETURNING id`,
			name, nullStr(cfg.Description), customerID, typ, nullStr(cfg.Password),
			nullStr(cfg.BackgroundColor), nullStr(cfg.TextColor), nullStr(cfg.BackgroundImageURL),
			nullStr(cfg.QRCodeKey), nullStr(cfg.BaseURL), nullStr(cfg.DefaultFilePath),
			nullInt(cfg.MainAppID), nullInt(cfg.ContentAppID), permissive, settings,
		).Scan(&id)
	} else {
		_, err = tx.ExecContext(ctx, `
			UPDATE configurations SET
				name=$1, description=$2, type=$3, password=$4,
				backgroundcolor=$5, textcolor=$6, backgroundimageurl=$7,
				qrcodekey=$8, baseurl=$9, defaultfilepath=$10,
				mainappid=$11, contentappid=$12, permissive=$13, settingsjson=$14
			WHERE id=$15 AND customerid=$16`,
			name, nullStr(cfg.Description), typ, nullStr(cfg.Password),
			nullStr(cfg.BackgroundColor), nullStr(cfg.TextColor), nullStr(cfg.BackgroundImageURL),
			nullStr(cfg.QRCodeKey), nullStr(cfg.BaseURL), nullStr(cfg.DefaultFilePath),
			nullInt(cfg.MainAppID), nullInt(cfg.ContentAppID), permissive, settings,
			id, customerID,
		)
	}
	if err != nil {
		return 0, err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM configurationapplicationparameters WHERE configurationid=$1`, id); err != nil {
		return 0, err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM configurationapplications WHERE configurationid=$1`, id); err != nil {
		return 0, err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM configurationfiles WHERE configurationid=$1`, id); err != nil {
		return 0, err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM configurationapplicationsettings WHERE configurationid=$1`, id); err != nil {
		return 0, err
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
		bottom := false
		if app.Bottom != nil {
			bottom = *app.Bottom
		}
		remove := false
		if app.Remove != nil {
			remove = *app.Remove
		}
		longTap := false
		if app.LongTap != nil {
			longTap = *app.LongTap
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO configurationapplications (
				configurationid, applicationid, applicationversionid,
				action, showicon, screenorder, keycode, bottom, remove, longtap
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
			id, appID, verID, action, showIcon,
			nullInt(app.ScreenOrder), nullInt(app.KeyCode), bottom, remove, longTap,
		)
		if err != nil {
			return 0, err
		}
		if app.SkipVersionCheck != nil {
			if err := r.upsertConfigAppParam(ctx, tx, id, appID, *app.SkipVersionCheck); err != nil {
				return 0, err
			}
		}
	}
	for _, f := range cfg.Files {
		if f.Remove != nil && *f.Remove {
			continue
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO configurationfiles (configurationid, path, externalurl, url, remove)
			VALUES ($1,$2,$3,$4,$5)`,
			id, nullStr(f.Path), nullStr(f.ExternalURL), nullStr(f.URL), false,
		)
		if err != nil {
			return 0, err
		}
	}
	for _, s := range cfg.ApplicationSettings {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO configurationapplicationsettings (configurationid, applicationid, name, type, value)
			VALUES ($1,$2,$3,$4,$5)`,
			id, nullInt(s.ApplicationID), nullStr(s.Name), nullStr(s.Type), nullStr(s.Value),
		)
		if err != nil {
			return 0, err
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ConfigRepository) upsertConfigAppParam(ctx context.Context, tx *sql.Tx, configurationID, applicationID int, skip bool) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO configurationapplicationparameters (configurationid, applicationid, skipversioncheck)
		VALUES ($1, $2, $3)
		ON CONFLICT (configurationid, applicationid) DO UPDATE SET skipversioncheck = $3`,
		configurationID, applicationID, skip)
	return err
}

func (r *ConfigRepository) Delete(ctx context.Context, customerID, id int) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM configurations WHERE id=$1 AND customerid=$2`, id, customerID)
	return err
}

func (r *ConfigRepository) Copy(ctx context.Context, customerID int, req domain.CopyRequest) (int, error) {
	src, err := r.GetByID(ctx, customerID, req.ID)
	if err != nil || src == nil {
		return 0, err
	}
	name := strings.TrimSpace(req.Name)
	src.ID = nil
	src.Name = &name
	src.QRCodeKey = nil
	if req.Description != nil {
		src.Description = req.Description
	}
	domain.EnsureQRCodeKey(src, nil)
	return r.Insert(ctx, customerID, *src)
}

func (r *ConfigRepository) ListAllApplicationsForPicker(ctx context.Context, customerID int) ([]domain.ConfigurationApplication, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT a.id, a.name, a.pkg, a.type,
		       lv.id AS latest_version,
		       lv.version AS latest_version_text
		FROM applications a
		LEFT JOIN LATERAL (
			SELECT av.id, av.version
			FROM applicationversions av
			WHERE av.applicationid = a.id
			ORDER BY av.versioncode DESC, av.id DESC
			LIMIT 1
		) lv ON TRUE
		WHERE a.customerid = $1 OR a.common = TRUE
		ORDER BY lower(a.name)`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanConfigApps(rows)
}

func (r *ConfigRepository) ListConfigurationApplications(ctx context.Context, customerID, configurationID int) ([]domain.ConfigurationApplication, error) {
	_ = customerID
	rows, err := r.db.QueryContext(ctx, `
		SELECT a.id, a.name, a.pkg, a.type,
		       ca.applicationversionid, ca.action, ca.showicon,
		       ca.screenorder, ca.keycode, ca.bottom,
		       ca.remove, ca.longtap,
		       COALESCE(cap.skipversioncheck, false),
		       av.version, av.versioncode, av.url,
		       (SELECT av2.id FROM applicationversions av2
		        WHERE av2.applicationid = a.id ORDER BY av2.versioncode DESC, av2.id DESC LIMIT 1) AS latest_version
		FROM configurationapplications ca
		JOIN applications a ON a.id = ca.applicationid
		LEFT JOIN applicationversions av ON av.id = ca.applicationversionid
		LEFT JOIN configurationapplicationparameters cap
		       ON cap.configurationid = ca.configurationid AND cap.applicationid = ca.applicationid
		WHERE ca.configurationid = $1
		ORDER BY lower(a.name)`, configurationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.ConfigurationApplication, 0)
	for rows.Next() {
		var app domain.ConfigurationApplication
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

func (r *ConfigRepository) UpgradeApplication(ctx context.Context, customerID int, req domain.UpgradeApplicationRequest) error {
	var latest int
	err := r.db.QueryRowContext(ctx, `
		SELECT av.id FROM applicationversions av
		JOIN applications a ON a.id = av.applicationid
		WHERE av.applicationid = $1 AND (a.customerid = $2 OR a.common = TRUE)
		ORDER BY av.versioncode DESC, av.id DESC LIMIT 1`,
		req.ApplicationID, customerID).Scan(&latest)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		UPDATE configurationapplications
		SET applicationversionid = $1
		WHERE configurationid = $2 AND applicationid = $3`,
		latest, req.ConfigurationID, req.ApplicationID)
	return err
}

func (r *ConfigRepository) loadFiles(ctx context.Context, configurationID int) ([]domain.ConfigurationFile, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, path, externalurl, url, remove
		FROM configurationfiles WHERE configurationid = $1`, configurationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ConfigurationFile
	for rows.Next() {
		var f domain.ConfigurationFile
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

func (r *ConfigRepository) loadAppSettings(ctx context.Context, configurationID int) ([]domain.ConfigurationApplicationSetting, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.applicationid, a.name, s.name, s.type, s.value
		FROM configurationapplicationsettings s
		LEFT JOIN applications a ON a.id = s.applicationid
		WHERE s.configurationid = $1`, configurationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ConfigurationApplicationSetting
	for rows.Next() {
		var s domain.ConfigurationApplicationSetting
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

func scanConfigApps(rows *sql.Rows) ([]domain.ConfigurationApplication, error) {
	out := make([]domain.ConfigurationApplication, 0)
	for rows.Next() {
		var app domain.ConfigurationApplication
		var name, pkg, typ, verText sql.NullString
		var latest sql.NullInt64
		if err := rows.Scan(&app.ID, &name, &pkg, &typ, &latest, &verText); err != nil {
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
		if latest.Valid {
			l := int(latest.Int64)
			app.LatestVersion = &l
		}
		if verText.Valid {
			app.Version = &verText.String
		}
		app.ApplicationID = app.ID
		out = append(out, app)
	}
	return out, rows.Err()
}

func nullStr(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}

func nullInt(n *int) any {
	if n == nil || *n == 0 {
		return nil
	}
	return *n
}
