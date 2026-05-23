package postgres

import (
	"context"
	"database/sql"

	appdomain "github.com/gis-mdm/server-backend-go/internal/modules/applications/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/files/port"
)

// AppLookup implements port.ApplicationLookup.
type AppLookup struct {
	db *sql.DB
}

func NewAppLookup(db *sql.DB) *AppLookup {
	return &AppLookup{db: db}
}

var _ port.ApplicationLookup = (*AppLookup)(nil)

func (a *AppLookup) SearchByURL(ctx context.Context, customerID int, url string) ([]appdomain.Application, error) {
	rows, err := a.db.QueryContext(ctx, `
		SELECT a.id, a.name, a.pkg, a.type, a.url, a.customerid
		FROM applications a
		WHERE (a.customerid = $1 OR a.common = TRUE)
		  AND (a.url = $2 OR EXISTS (
		      SELECT 1 FROM applicationversions av
		      WHERE av.applicationid = a.id AND (av.url = $2 OR av.urlarmeabi = $2 OR av.urlarm64 = $2)
		  ))`, customerID, url)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []appdomain.Application
	for rows.Next() {
		var app appdomain.Application
		var urlVal sql.NullString
		if err := rows.Scan(&app.ID, &app.Name, &app.Pkg, &app.Type, &urlVal, &app.CustomerID); err != nil {
			return nil, err
		}
		if urlVal.Valid {
			app.URL = &urlVal.String
		}
		out = append(out, app)
	}
	return out, rows.Err()
}

func (a *AppLookup) FindVersionByPkgCode(ctx context.Context, customerID int, pkg string, versionCode int) (*appdomain.ApplicationVersion, error) {
	return a.findVersion(ctx, customerID, pkg, "versioncode = $3", versionCode)
}

func (a *AppLookup) FindVersionByPkgVersion(ctx context.Context, customerID int, pkg, version string) (*appdomain.ApplicationVersion, error) {
	return a.findVersion(ctx, customerID, pkg, "av.version = $3", version)
}

func (a *AppLookup) findVersion(ctx context.Context, customerID int, pkg, clause string, arg any) (*appdomain.ApplicationVersion, error) {
	q := `
		SELECT av.id, av.applicationid, av.version, av.versioncode, av.url, av.urlarmeabi, av.urlarm64, av.split
		FROM applicationversions av
		INNER JOIN applications a ON a.id = av.applicationid
		WHERE (a.customerid = $1 OR a.common = TRUE) AND lower(a.pkg) = lower($2) AND ` + clause + `
		LIMIT 1`
	row := a.db.QueryRowContext(ctx, q, customerID, pkg, arg)
	var v appdomain.ApplicationVersion
	var url, armeabi, arm64 sql.NullString
	var split bool
	if err := row.Scan(&v.ID, &v.ApplicationID, &v.Version, &v.VersionCode, &url, &armeabi, &arm64, &split); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if url.Valid {
		v.URL = &url.String
	}
	if armeabi.Valid {
		v.URLArmeabi = &armeabi.String
	}
	if arm64.Valid {
		v.URLArm64 = &arm64.String
	}
	v.Split = &split
	return &v, nil
}

func (a *AppLookup) FindAppsByPkg(ctx context.Context, customerID int, pkg string) ([]appdomain.Application, error) {
	rows, err := a.db.QueryContext(ctx, `
		SELECT a.id, a.name, a.pkg, a.showicon, a.usekiosk, a.runafterinstall, a.runatboot, a.system
		FROM applications a
		WHERE (a.customerid = $1 OR a.common = TRUE) AND lower(a.pkg) = lower($2)
		LIMIT 1`, customerID, pkg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []appdomain.Application
	for rows.Next() {
		var app appdomain.Application
		var showIcon, useKiosk, runAfter, runAtBoot, system bool
		if err := rows.Scan(&app.ID, &app.Name, &app.Pkg, &showIcon, &useKiosk, &runAfter, &runAtBoot, &system); err != nil {
			return nil, err
		}
		app.ShowIcon = &showIcon
		app.UseKiosk = &useKiosk
		app.RunAfterInstall = &runAfter
		app.RunAtBoot = &runAtBoot
		app.System = &system
		out = append(out, app)
	}
	return out, rows.Err()
}
