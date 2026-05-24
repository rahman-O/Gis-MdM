package application

import (
	"context"
	"database/sql"
)

// ProfileVersionIDForEnrollmentRoute returns the published profile version bound to a route (FR-008a stub).
func ProfileVersionIDForEnrollmentRoute(ctx context.Context, db *sql.DB, enrollmentRouteID int64) (int, error) {
	if enrollmentRouteID <= 0 {
		return 0, sql.ErrNoRows
	}
	var profileVersionID sql.NullInt64
	err := db.QueryRowContext(ctx, `
		SELECT profile_version_id FROM enrollment_routes WHERE id = $1`, enrollmentRouteID).
		Scan(&profileVersionID)
	if err != nil {
		return 0, err
	}
	if !profileVersionID.Valid || profileVersionID.Int64 <= 0 {
		return 0, sql.ErrNoRows
	}
	return int(profileVersionID.Int64), nil
}
