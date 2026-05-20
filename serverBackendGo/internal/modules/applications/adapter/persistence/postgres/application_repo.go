package postgres

import (
	"context"
	"database/sql"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/applications/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/applications/port"
)

// ApplicationRepository implements port.ApplicationRepository.
type ApplicationRepository struct {
	db *sql.DB
}

func NewApplicationRepository(db *sql.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

var _ port.ApplicationRepository = (*ApplicationRepository)(nil)

const appSelect = `
	SELECT a.id, a.name, a.pkg, a.type, a.showicon, a.system, a.url, a.intent,
	       a.customerid, a.common, a.runafterinstall, a.runatboot, a.usekiosk, a.skipversion,
	       a.icontext, a.iconid,
	       lv.id AS latest_version_id, lv.version AS latest_version_text, lv.versioncode,
	       lv.url AS version_url
	FROM applications a
	LEFT JOIN LATERAL (
		SELECT av.id, av.version, av.versioncode, av.url
		FROM applicationversions av
		WHERE av.applicationid = a.id
		ORDER BY av.versioncode DESC, av.id DESC
		LIMIT 1
	) lv ON TRUE`

func (r *ApplicationRepository) Search(ctx context.Context, customerID int) ([]domain.Application, error) {
	return r.search(ctx, customerID, "")
}

func (r *ApplicationRepository) SearchByValue(ctx context.Context, customerID int, value string) ([]domain.Application, error) {
	return r.search(ctx, customerID, strings.TrimSpace(value))
}

func (r *ApplicationRepository) search(ctx context.Context, customerID int, filter string) ([]domain.Application, error) {
	q := appSelect + ` WHERE (a.customerid = $1 OR a.common = TRUE)`
	args := []any{customerID}
	if filter != "" {
		q += ` AND (a.name ILIKE $2 OR a.pkg ILIKE $2)`
		args = append(args, "%"+filter+"%")
	}
	q += ` ORDER BY lower(a.name)`
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanApps(rows)
}

func (r *ApplicationRepository) GetByID(ctx context.Context, customerID, id int) (*domain.Application, error) {
	row := r.db.QueryRowContext(ctx, appSelect+` WHERE a.id = $1 AND (a.customerid = $2 OR a.common = TRUE)`, id, customerID)
	apps, err := scanAppsRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(apps) == 0 {
		return nil, nil
	}
	return &apps[0], nil
}

func (r *ApplicationRepository) ListVersions(ctx context.Context, customerID, applicationID int) ([]domain.ApplicationVersion, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT av.id, av.applicationid, av.version, av.versioncode, av.url,
		       av.urlarmeabi, av.urlarm64, av.filepath, av.split, av.arch,
		       av.action, av.showicon, av.screenorder, av.keycode, av.bottom, av.autoupdate
		FROM applicationversions av
		JOIN applications a ON a.id = av.applicationid
		WHERE av.applicationid = $1 AND (a.customerid = $2 OR a.common = TRUE)
		ORDER BY av.versioncode DESC, av.id DESC`, applicationID, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ApplicationVersion
	for rows.Next() {
		var v domain.ApplicationVersion
		var id, appID, verCode sql.NullInt64
		var version, url, urlA, url64, fp, arch sql.NullString
		var split, showIcon, bottom, autoUpdate sql.NullBool
		var action, screenOrder, keyCode sql.NullInt64
		if err := rows.Scan(
			&id, &appID, &version, &verCode, &url, &urlA, &url64, &fp, &split, &arch,
			&action, &showIcon, &screenOrder, &keyCode, &bottom, &autoUpdate,
		); err != nil {
			return nil, err
		}
		if id.Valid {
			i := int(id.Int64)
			v.ID = &i
		}
		if appID.Valid {
			i := int(appID.Int64)
			v.ApplicationID = &i
		}
		if version.Valid {
			v.Version = &version.String
		}
		if verCode.Valid {
			i := int(verCode.Int64)
			v.VersionCode = &i
		}
		if url.Valid {
			v.URL = &url.String
		}
		if urlA.Valid {
			v.URLArmeabi = &urlA.String
		}
		if url64.Valid {
			v.URLArm64 = &url64.String
		}
		if fp.Valid {
			v.FilePath = &fp.String
		}
		if split.Valid {
			v.Split = &split.Bool
		}
		if arch.Valid {
			v.Arch = &arch.String
		}
		if action.Valid {
			i := int(action.Int64)
			v.Action = &i
		}
		if showIcon.Valid {
			v.ShowIcon = &showIcon.Bool
		}
		if screenOrder.Valid {
			i := int(screenOrder.Int64)
			v.ScreenOrder = &i
		}
		if keyCode.Valid {
			i := int(keyCode.Int64)
			v.KeyCode = &i
		}
		if bottom.Valid {
			v.Bottom = &bottom.Bool
		}
		if autoUpdate.Valid {
			v.AutoUpdate = &autoUpdate.Bool
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (r *ApplicationRepository) SaveAndroid(ctx context.Context, customerID int, app domain.Application) (*domain.Application, error) {
	t := "app"
	app.Type = &t
	return r.saveApp(ctx, customerID, app)
}

func (r *ApplicationRepository) SaveWeb(ctx context.Context, customerID int, app domain.Application) (*domain.Application, error) {
	t := "web"
	app.Type = &t
	return r.saveApp(ctx, customerID, app)
}

func (r *ApplicationRepository) saveApp(ctx context.Context, customerID int, app domain.Application) (*domain.Application, error) {
	pkg := deref(app.Pkg)
	name := deref(app.Name)
	showIcon := true
	if app.ShowIcon != nil {
		showIcon = *app.ShowIcon
	}
	system := false
	if app.System != nil {
		system = *app.System
	}
	runAfter := false
	if app.RunAfterInstall != nil {
		runAfter = *app.RunAfterInstall
	}
	runBoot := false
	if app.RunAtBoot != nil {
		runBoot = *app.RunAtBoot
	}
	useKiosk := false
	if app.UseKiosk != nil {
		useKiosk = *app.UseKiosk
	}
	skip := false
	if app.SkipVersion != nil {
		skip = *app.SkipVersion
	}
	id := 0
	if app.ID != nil {
		id = *app.ID
	}
	var err error
	if id == 0 {
		err = r.db.QueryRowContext(ctx, `
			INSERT INTO applications (
				pkg, name, customerid, type, showicon, system, url, intent,
				runafterinstall, runatboot, usekiosk, skipversion, icontext, iconid
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
			RETURNING id`,
			pkg, name, customerID, deref(app.Type), showIcon, system,
			nullStr(app.URL), nullStr(app.Intent), runAfter, runBoot, useKiosk, skip,
			nullStr(app.IconText), nullInt(app.IconID),
		).Scan(&id)
	} else {
		_, err = r.db.ExecContext(ctx, `
			UPDATE applications SET
				pkg=$1, name=$2, type=$3, showicon=$4, system=$5, url=$6, intent=$7,
				runafterinstall=$8, runatboot=$9, usekiosk=$10, skipversion=$11,
				icontext=$12, iconid=$13
			WHERE id=$14 AND customerid=$15`,
			pkg, name, deref(app.Type), showIcon, system,
			nullStr(app.URL), nullStr(app.Intent), runAfter, runBoot, useKiosk, skip,
			nullStr(app.IconText), nullInt(app.IconID), id, customerID,
		)
	}
	if err != nil {
		return nil, err
	}
	if app.Version != nil || app.VersionCode != nil {
		ver := domain.ApplicationVersion{
			ApplicationID: &id,
			Version:       app.Version,
			VersionCode:   app.VersionCode,
			URL:           app.URL,
		}
		if app.ID != nil && app.UsedVersionID != nil {
			ver.ID = app.UsedVersionID
		}
		_, _ = r.SaveVersion(ctx, customerID, ver)
	}
	return r.GetByID(ctx, customerID, id)
}

func (r *ApplicationRepository) SaveVersion(ctx context.Context, customerID int, ver domain.ApplicationVersion) (*domain.ApplicationVersion, error) {
	appID := 0
	if ver.ApplicationID != nil {
		appID = *ver.ApplicationID
	}
	verCode := 0
	if ver.VersionCode != nil {
		verCode = *ver.VersionCode
	}
	id := 0
	if ver.ID != nil {
		id = *ver.ID
	}
	var err error
	if id == 0 {
		var ok int
		if err := r.db.QueryRowContext(ctx, `
			SELECT a.id FROM applications a
			WHERE a.id = $1 AND (a.customerid = $2 OR a.common = TRUE)`, appID, customerID).Scan(&ok); err != nil {
			return nil, err
		}
		err = r.db.QueryRowContext(ctx, `
			INSERT INTO applicationversions (
				applicationid, version, versioncode, url, urlarmeabi, urlarm64,
				filepath, split, arch, action, showicon, screenorder, keycode, bottom, autoupdate
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
			RETURNING id`,
			appID, nullStr(ver.Version), verCode, nullStr(ver.URL), nullStr(ver.URLArmeabi),
			nullStr(ver.URLArm64), nullStr(ver.FilePath), boolOrFalse(ver.Split), nullStr(ver.Arch),
			nullInt(ver.Action), boolOrFalse(ver.ShowIcon), nullInt(ver.ScreenOrder),
			nullInt(ver.KeyCode), boolOrFalse(ver.Bottom), boolOrFalse(ver.AutoUpdate),
		).Scan(&id)
	} else {
		_, err = r.db.ExecContext(ctx, `
			UPDATE applicationversions av SET
				version=$1, versioncode=$2, url=$3, urlarmeabi=$4, urlarm64=$5,
				filepath=$6, split=$7, arch=$8, action=$9, showicon=$10,
				screenorder=$11, keycode=$12, bottom=$13, autoupdate=$14
			FROM applications a
			WHERE av.id=$15 AND av.applicationid=a.id
			  AND (a.customerid=$16 OR a.common=TRUE)`,
			nullStr(ver.Version), verCode, nullStr(ver.URL), nullStr(ver.URLArmeabi),
			nullStr(ver.URLArm64), nullStr(ver.FilePath), boolOrFalse(ver.Split), nullStr(ver.Arch),
			nullInt(ver.Action), boolOrFalse(ver.ShowIcon), nullInt(ver.ScreenOrder),
			nullInt(ver.KeyCode), boolOrFalse(ver.Bottom), boolOrFalse(ver.AutoUpdate),
			id, customerID,
		)
	}
	if err != nil {
		return nil, err
	}
	_, _ = r.db.ExecContext(ctx, `UPDATE applications SET latestversion=$1 WHERE id=$2`, id, appID)
	out := ver
	out.ID = &id
	return &out, nil
}

func (r *ApplicationRepository) DeleteApp(ctx context.Context, customerID, id int) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM applications WHERE id=$1 AND customerid=$2`, id, customerID)
	return err
}

func (r *ApplicationRepository) DeleteVersion(ctx context.Context, customerID, versionID int) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM applicationversions av
		USING applications a
		WHERE av.id=$1 AND av.applicationid=a.id AND (a.customerid=$2 OR a.common=TRUE)`,
		versionID, customerID)
	return err
}

func (r *ApplicationRepository) ValidatePkg(ctx context.Context, customerID int, req domain.ValidatePkgRequest) ([]domain.Application, error) {
	rows, err := r.db.QueryContext(ctx, appSelect+`
		WHERE (a.customerid = $1 OR a.common = TRUE) AND lower(a.pkg) = lower($2)
		  AND ($3 = 0 OR a.id <> $3)`,
		customerID, req.Pkg, idOrZero(req.ID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanApps(rows)
}

func (r *ApplicationRepository) GetAppConfigurations(ctx context.Context, customerID, applicationID int) ([]domain.ApplicationConfigurationLink, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT ca.id, ca.configurationid, c.name, ca.action
		FROM configurationapplications ca
		JOIN configurations c ON c.id = ca.configurationid
		JOIN applications a ON a.id = ca.applicationid
		WHERE ca.applicationid = $1 AND c.customerid = $2`,
		applicationID, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ApplicationConfigurationLink
	for rows.Next() {
		var l domain.ApplicationConfigurationLink
		var id int
		var name sql.NullString
		var action int
		if err := rows.Scan(&id, &l.ConfigurationID, &name, &action); err != nil {
			return nil, err
		}
		l.ID = &id
		aid := applicationID
		l.ApplicationID = &aid
		if name.Valid {
			l.Name = &name.String
		}
		l.Action = action
		sel := true
		l.Selected = &sel
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *ApplicationRepository) UpdateAppConfigurations(ctx context.Context, customerID int, req domain.LinkConfigurationsToAppRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `
		DELETE FROM configurationapplications
		WHERE applicationid = $1 AND configurationid IN (
			SELECT id FROM configurations WHERE customerid = $2)`,
		req.ApplicationID, customerID)
	if err != nil {
		return err
	}
	for _, l := range req.Configurations {
		if l.Selected != nil && !*l.Selected {
			continue
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO configurationapplications (configurationid, applicationid, action, showicon)
			VALUES ($1,$2,$3,TRUE)
			ON CONFLICT (configurationid, applicationid) DO UPDATE SET action = EXCLUDED.action`,
			l.ConfigurationID, req.ApplicationID, l.Action)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *ApplicationRepository) GetVersionConfigurations(ctx context.Context, customerID, versionID int) ([]domain.ApplicationVersionConfigurationLink, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT ca.configurationid, c.name, ca.action
		FROM configurationapplications ca
		JOIN configurations c ON c.id = ca.configurationid
		JOIN applicationversions av ON av.applicationid = ca.applicationid
		WHERE av.id = $1 AND c.customerid = $2 AND ca.applicationversionid = $1`,
		versionID, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ApplicationVersionConfigurationLink
	for rows.Next() {
		var l domain.ApplicationVersionConfigurationLink
		var name sql.NullString
		if err := rows.Scan(&l.ConfigurationID, &name, &l.Action); err != nil {
			return nil, err
		}
		if name.Valid {
			l.Name = &name.String
		}
		sel := true
		l.Selected = &sel
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *ApplicationRepository) UpdateVersionConfigurations(ctx context.Context, customerID int, req domain.LinkConfigurationsToAppVersionRequest) error {
	var appID int
	err := r.db.QueryRowContext(ctx, `
		SELECT applicationid FROM applicationversions av
		JOIN applications a ON a.id = av.applicationid
		WHERE av.id = $1 AND (a.customerid = $2 OR a.common = TRUE)`,
		req.ApplicationVersionID, customerID).Scan(&appID)
	if err != nil {
		return err
	}
	link := domain.LinkConfigurationsToAppRequest{ApplicationID: appID}
	for _, l := range req.Configurations {
		link.Configurations = append(link.Configurations, domain.ApplicationConfigurationLink{
			ConfigurationID: l.ConfigurationID,
			Action:          l.Action,
			Selected:        l.Selected,
		})
	}
	return r.UpdateAppConfigurations(ctx, customerID, link)
}

func (r *ApplicationRepository) AdminSearch(ctx context.Context, value string) ([]domain.Application, error) {
	q := appSelect + ` WHERE a.common = TRUE`
	args := []any{}
	if strings.TrimSpace(value) != "" {
		q += ` AND (a.name ILIKE $1 OR a.pkg ILIKE $1)`
		args = append(args, "%"+strings.TrimSpace(value)+"%")
	}
	q += ` ORDER BY lower(a.name)`
	var rows *sql.Rows
	var err error
	if len(args) == 0 {
		rows, err = r.db.QueryContext(ctx, q)
	} else {
		rows, err = r.db.QueryContext(ctx, q, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanApps(rows)
}

func (r *ApplicationRepository) TurnIntoCommon(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE applications SET common = TRUE, customerid = NULL WHERE id = $1`, id)
	return err
}

func scanApps(rows *sql.Rows) ([]domain.Application, error) {
	var out []domain.Application
	for rows.Next() {
		app, err := scanAppRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, app)
	}
	return out, rows.Err()
}

func scanAppsRows(row *sql.Row) ([]domain.Application, error) {
	app, err := scanAppFromRow(row)
	if err != nil {
		return nil, err
	}
	return []domain.Application{app}, nil
}

func scanAppRow(rows *sql.Rows) (domain.Application, error) {
	var app domain.Application
	var id int
	var name, pkg, typ, url, intent, iconText sql.NullString
	var showIcon, system, common, runAfter, runBoot, useKiosk, skip sql.NullBool
	var customerID, iconID sql.NullInt64
	var latestID, verCode sql.NullInt64
	var latestText, versionURL sql.NullString
	err := rows.Scan(
		&id, &name, &pkg, &typ, &showIcon, &system, &url, &intent,
		&customerID, &common, &runAfter, &runBoot, &useKiosk, &skip,
		&iconText, &iconID,
		&latestID, &latestText, &verCode, &versionURL,
	)
	if err != nil {
		return app, err
	}
	return fillApp(id, name, pkg, typ, url, intent, iconText,
		showIcon, system, customerID, common, runAfter, runBoot, useKiosk, skip,
		iconID, latestID, verCode, latestText, versionURL), nil
}

func scanAppFromRow(row *sql.Row) (domain.Application, error) {
	var app domain.Application
	var id int
	var name, pkg, typ, url, intent, iconText sql.NullString
	var showIcon, system, common, runAfter, runBoot, useKiosk, skip sql.NullBool
	var customerID, iconID sql.NullInt64
	var latestID, verCode sql.NullInt64
	var latestText, versionURL sql.NullString
	err := row.Scan(
		&id, &name, &pkg, &typ, &showIcon, &system, &url, &intent,
		&customerID, &common, &runAfter, &runBoot, &useKiosk, &skip,
		&iconText, &iconID,
		&latestID, &latestText, &verCode, &versionURL,
	)
	if err != nil {
		return app, err
	}
	return fillApp(id, name, pkg, typ, url, intent, iconText,
		showIcon, system, customerID, common, runAfter, runBoot, useKiosk, skip,
		iconID, latestID, verCode, latestText, versionURL), nil
}

func fillApp(id int, name, pkg, typ, url, intent, iconText sql.NullString,
	showIcon, system sql.NullBool,
	customerID sql.NullInt64, common sql.NullBool,
	runAfter, runBoot, useKiosk, skip sql.NullBool,
	iconID, latestID, verCode sql.NullInt64,
	latestText, versionURL sql.NullString) domain.Application {
	app := domain.Application{ID: &id}
	if name.Valid {
		app.Name = &name.String
	}
	if pkg.Valid {
		app.Pkg = &pkg.String
	}
	if typ.Valid {
		app.Type = &typ.String
	}
	if showIcon.Valid {
		app.ShowIcon = &showIcon.Bool
	}
	if system.Valid {
		app.System = &system.Bool
	}
	if url.Valid {
		app.URL = &url.String
	}
	if intent.Valid {
		app.Intent = &intent.String
	}
	if customerID.Valid {
		c := int(customerID.Int64)
		app.CustomerID = &c
	}
	if common.Valid {
		app.Common = &common.Bool
		app.CommonApplication = &common.Bool
	}
	if runAfter.Valid {
		app.RunAfterInstall = &runAfter.Bool
	}
	if runBoot.Valid {
		app.RunAtBoot = &runBoot.Bool
	}
	if useKiosk.Valid {
		app.UseKiosk = &useKiosk.Bool
	}
	if skip.Valid {
		app.SkipVersion = &skip.Bool
	}
	if iconText.Valid {
		app.IconText = &iconText.String
	}
	if iconID.Valid {
		i := int(iconID.Int64)
		app.IconID = &i
	}
	if latestID.Valid {
		l := int(latestID.Int64)
		app.LatestVersion = &l
		app.UsedVersionID = &l
	}
	if latestText.Valid {
		app.LatestVersionText = &latestText.String
		app.Version = &latestText.String
	}
	if verCode.Valid {
		v := int(verCode.Int64)
		app.VersionCode = &v
	}
	if versionURL.Valid && versionURL.String != "" {
		app.URL = &versionURL.String
	}
	return app
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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

func boolOrFalse(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func idOrZero(id *int) int {
	if id == nil {
		return 0
	}
	return *id
}
