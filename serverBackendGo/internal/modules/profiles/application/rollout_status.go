package application

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	profilepostgres "github.com/gis-mdm/server-backend-go/internal/modules/profiles/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/port"
	devdomain "github.com/gis-mdm/server-backend-go/internal/modules/devices/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// RolloutStatusService recomputes device rollout status from sync telemetry.
type RolloutStatusService struct {
	store port.RolloutStore
	db    *sql.DB
}

func NewRolloutStatusService(db *sql.DB) *RolloutStatusService {
	return &RolloutStatusService{store: profilepostgres.NewAssignmentRepository(db), db: db}
}

func (s *RolloutStatusService) ListDevices(ctx context.Context, p *platformauth.Principal, profileID int, q domain.RolloutDevicesQuery) (*domain.RolloutDevicesPage, error) {
	if err := requireConfigPerm(p); err != nil {
		return nil, err
	}
	items, total, err := s.store.ListRolloutDevices(ctx, customerID(p), profileID, q)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.DeviceRolloutRow{}
	}
	return &domain.RolloutDevicesPage{Items: items, TotalCount: total}, nil
}

func (s *RolloutStatusService) RecomputeProfile(ctx context.Context, p *platformauth.Principal, profileID int) error {
	if err := requireConfigPerm(p); err != nil {
		return err
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id FROM devices WHERE customerid = $1 AND (
			target_profile_version_id IN (SELECT profile_version_id FROM profile_tree_assignments WHERE profile_id = $2)
			OR applied_profile_version_id IN (SELECT id FROM profile_versions WHERE profile_id = $2)
		)`, customerID(p), profileID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return err
		}
		_ = RecomputeDeviceRollout(ctx, s.db, id)
	}
	return rows.Err()
}

// RecomputeDeviceRollout updates rollout columns for one device (called from sync/info).
func RecomputeDeviceRollout(ctx context.Context, db *sql.DB, deviceID int64) error {
	store := profilepostgres.NewAssignmentRepository(db)
	var targetID, appliedID sql.NullInt64
	var info sql.NullString
	err := db.QueryRowContext(ctx, `
		SELECT target_profile_version_id, applied_profile_version_id, info
		FROM devices WHERE id = $1`, deviceID).Scan(&targetID, &appliedID, &info)
	if err != nil {
		return err
	}
	if !targetID.Valid || targetID.Int64 <= 0 {
		return store.UpdateRolloutFields(ctx, deviceID, domain.RolloutPending, "No profile target assigned")
	}
	tgt := int(targetID.Int64)
	if !appliedID.Valid || appliedID.Int64 != int64(tgt) {
		return store.UpdateRolloutFields(ctx, deviceID, domain.RolloutPending, "Awaiting sync")
	}
	partial, reason := checkProfileApps(ctx, db, tgt, info.String)
	if partial {
		return store.UpdateRolloutFields(ctx, deviceID, domain.RolloutPartial, reason)
	}
	return store.UpdateRolloutFields(ctx, deviceID, domain.RolloutInstalled, "")
}

func checkProfileApps(ctx context.Context, db *sql.DB, versionID int, infoJSON string) (bool, string) {
	rows, err := db.QueryContext(ctx, `
		SELECT a.pkg FROM profile_version_applications pva
		JOIN applications a ON a.id = pva.applicationid
		WHERE pva.profile_version_id = $1 AND pva.remove = false`, versionID)
	if err != nil {
		return false, ""
	}
	defer rows.Close()
	want := make(map[string]struct{})
	for rows.Next() {
		var pkg string
		if err := rows.Scan(&pkg); err != nil {
			return false, ""
		}
		if pkg != "" {
			want[strings.ToLower(pkg)] = struct{}{}
		}
	}
	if len(want) == 0 {
		return false, ""
	}
	var info devdomain.DeviceInfoView
	if infoJSON != "" {
		_ = json.Unmarshal([]byte(infoJSON), &info)
	}
	installed := make(map[string]string)
	for _, app := range info.Applications {
		if app.Pkg == nil {
			continue
		}
		pkg := strings.ToLower(strings.TrimSpace(*app.Pkg))
		st := ""
		if app.Status != nil {
			st = strings.ToUpper(strings.TrimSpace(*app.Status))
		}
		installed[pkg] = st
	}
	var problems []string
	for pkg := range want {
		st, ok := installed[pkg]
		if !ok || st == "" {
			problems = append(problems, pkg+": not reported")
			continue
		}
		if strings.Contains(st, "FAIL") || st == "3" || strings.Contains(st, "MISMATCH") || st == "2" {
			problems = append(problems, pkg+": "+st)
		}
	}
	if len(problems) == 0 {
		return false, ""
	}
	reason := problems[0]
	if len(problems) > 1 {
		reason += fmt.Sprintf(" (+%d more)", len(problems)-1)
	}
	return true, reason
}
