package application

import (
	"context"
	"database/sql"

	profilepostgres "github.com/gis-mdm/server-backend-go/internal/modules/profiles/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
)

// ResolveEffectiveProfile picks tree assignment over enrollment route (018 R1).
func ResolveEffectiveProfile(ctx context.Context, db *sql.DB, deviceID int64) (domain.EffectiveProfileResolution, error) {
	out := domain.EffectiveProfileResolution{Source: "none", Enabled: true}
	store := profilepostgres.NewAssignmentRepository(db)
	dc, err := store.LoadDeviceContext(ctx, deviceID)
	if err != nil {
		return out, err
	}
	var treeID int64
	if dc.TreeNodeID.Valid {
		treeID = dc.TreeNodeID.Int64
	}
	if pid, vid, ok, err := store.ResolveTreeVersion(ctx, dc.CustomerID, treeID); err != nil {
		return out, err
	} else if ok {
		out.ProfileID = pid
		out.ProfileVersionID = vid
		out.Source = "tree"
		enabled, _ := store.IsProfileEnabled(ctx, dc.CustomerID, pid)
		out.Enabled = enabled
		return out, nil
	}
	if dc.EnrollmentRouteID.Valid && dc.EnrollmentRouteID.Int64 > 0 {
		var profileID, versionID int
		var enabled bool
		err := db.QueryRowContext(ctx, `
			SELECT p.id, er.profile_version_id, COALESCE(p.enabled, true)
			FROM enrollment_routes er
			JOIN profile_versions pv ON pv.id = er.profile_version_id
			JOIN profiles p ON p.id = pv.profile_id
			WHERE er.id = $1 AND er.customerid = $2`, dc.EnrollmentRouteID.Int64, dc.CustomerID).
			Scan(&profileID, &versionID, &enabled)
		if err == nil && versionID > 0 {
			out.ProfileID = profileID
			out.ProfileVersionID = versionID
			out.RouteID = dc.EnrollmentRouteID.Int64
			out.Source = "route"
			out.Enabled = enabled
		}
	}
	return out, nil
}
