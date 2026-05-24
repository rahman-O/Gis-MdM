package application

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/domain"
)

var ErrStableVersionMissing = errors.New("error.enrollment_route.stable_version_missing")

// ResolveBootstrapIntent picks application version id for mainappid column.
func ResolveBootstrapIntent(ctx context.Context, db *sql.DB, customerID, applicationID int, intent string, specificVersionID *int) (domain.ResolvedBootstrap, error) {
	var out domain.ResolvedBootstrap
	if applicationID <= 0 {
		return out, ErrMainAppRequired
	}
	if !domain.ValidBootstrapIntent(intent) {
		intent = domain.BootstrapIntentStable
	}

	var pkg, appName string
	err := db.QueryRowContext(ctx, `
		SELECT pkg, name FROM applications
		WHERE id = $1 AND (customerid IS NULL OR customerid = $2)`, applicationID, customerID).
		Scan(&pkg, &appName)
	if err == sql.ErrNoRows {
		return out, ErrMainAppRequired
	}
	if err != nil {
		return out, err
	}
	out.ApplicationID = applicationID
	out.Package = pkg

	switch intent {
	case domain.BootstrapIntentSpecific:
		if specificVersionID == nil || *specificVersionID <= 0 {
			return out, ErrMainAppRequired
		}
		return loadVersion(ctx, db, applicationID, *specificVersionID, out)
	case domain.BootstrapIntentLatest:
		return loadLatestVersion(ctx, db, applicationID, out)
	default:
		return loadStableVersion(ctx, db, applicationID, out)
	}
}

func loadVersion(ctx context.Context, db *sql.DB, applicationID, versionID int, out domain.ResolvedBootstrap) (domain.ResolvedBootstrap, error) {
	err := db.QueryRowContext(ctx, `
		SELECT id, COALESCE(version, ''), versioncode
		FROM applicationversions
		WHERE id = $1 AND applicationid = $2`, versionID, applicationID).
		Scan(&out.VersionID, &out.VersionLabel, &out.VersionCode)
	if err == sql.ErrNoRows {
		return out, ErrMainAppRequired
	}
	if err != nil {
		return out, err
	}
	return out, nil
}

func loadStableVersion(ctx context.Context, db *sql.DB, applicationID int, out domain.ResolvedBootstrap) (domain.ResolvedBootstrap, error) {
	err := db.QueryRowContext(ctx, `
		SELECT id, COALESCE(version, ''), versioncode
		FROM applicationversions
		WHERE applicationid = $1 AND is_recommended = TRUE
		ORDER BY versioncode DESC
		LIMIT 1`, applicationID).
		Scan(&out.VersionID, &out.VersionLabel, &out.VersionCode)
	if err == sql.ErrNoRows {
		return out, ErrStableVersionMissing
	}
	if err != nil {
		return out, err
	}
	return out, nil
}

func loadLatestVersion(ctx context.Context, db *sql.DB, applicationID int, out domain.ResolvedBootstrap) (domain.ResolvedBootstrap, error) {
	err := db.QueryRowContext(ctx, `
		SELECT id, COALESCE(version, ''), versioncode
		FROM applicationversions
		WHERE applicationid = $1
		  AND (url IS NOT NULL AND trim(url) <> '')
		ORDER BY versioncode DESC
		LIMIT 1`, applicationID).
		Scan(&out.VersionID, &out.VersionLabel, &out.VersionCode)
	if err == sql.ErrNoRows {
		return out, ErrMainAppRequired
	}
	if err != nil {
		return out, err
	}
	return out, nil
}
