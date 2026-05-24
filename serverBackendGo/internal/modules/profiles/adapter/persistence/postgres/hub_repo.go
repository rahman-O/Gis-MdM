package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
)

// HubRepository loads aggregated hub metrics (019).
type HubRepository struct {
	db *sql.DB
}

func NewHubRepository(db *sql.DB) *HubRepository {
	return &HubRepository{db: db}
}

// ListRowMetrics per profile for list enrichment.
type ListRowMetrics struct {
	ProfileID           int
	AssignmentCount     int
	RolloutFailureCount int
	HasPublished        bool
	HasUnpublishedDraft bool
	StalePublish        bool
}

func (r *HubRepository) ListRowMetrics(ctx context.Context, customerID int, staleDays int) (map[int]ListRowMetrics, error) {
	staleInterval := fmt.Sprintf("%d days", staleDays)
	rows, err := r.db.QueryContext(ctx, `
		SELECT p.id,
		       COALESCE((SELECT COUNT(*)::int FROM profile_tree_assignments a
		                 WHERE a.profile_id = p.id AND a.customerid = p.customerid), 0),
		       COALESCE((SELECT COUNT(*)::int FROM devices d
		                 WHERE d.customerid = p.customerid
		                   AND d.profile_rollout_status = 'failed'
		                   AND d.target_profile_version_id IN (
		                       SELECT pv.id FROM profile_versions pv WHERE pv.profile_id = p.id)), 0),
		       (p.published_version_id IS NOT NULL AND p.published_version_id > 0),
		       (p.draft_version_id IS NOT NULL AND EXISTS (
		           SELECT 1 FROM profile_versions dv
		           WHERE dv.id = p.draft_version_id AND dv.status = 'draft')),
		       (pub.published_at IS NOT NULL AND pub.published_at < NOW() - $2::interval)
		FROM profiles p
		LEFT JOIN profile_versions pub ON pub.id = p.published_version_id
		WHERE p.customerid = $1`, customerID, staleInterval)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make(map[int]ListRowMetrics)
	for rows.Next() {
		var m ListRowMetrics
		if err := rows.Scan(&m.ProfileID, &m.AssignmentCount, &m.RolloutFailureCount,
			&m.HasPublished, &m.HasUnpublishedDraft, &m.StalePublish); err != nil {
			return nil, err
		}
		out[m.ProfileID] = m
	}
	return out, rows.Err()
}

func (r *HubRepository) RolloutSnapshot(ctx context.Context, customerID, profileID int) (domain.ProfileRolloutSnap, error) {
	var snap domain.ProfileRolloutSnap
	err := r.db.QueryRowContext(ctx, `
		SELECT
		  COALESCE(SUM(CASE WHEN d.profile_rollout_status = 'pending' THEN 1 ELSE 0 END), 0)::int,
		  COALESCE(SUM(CASE WHEN d.profile_rollout_status = 'installed' THEN 1 ELSE 0 END), 0)::int,
		  COALESCE(SUM(CASE WHEN d.profile_rollout_status = 'partial' THEN 1 ELSE 0 END), 0)::int,
		  COALESCE(SUM(CASE WHEN d.profile_rollout_status = 'failed' THEN 1 ELSE 0 END), 0)::int,
		  COUNT(*)::int
		FROM devices d
		WHERE d.customerid = $1
		  AND d.target_profile_version_id IN (
		      SELECT pv.id FROM profile_versions pv WHERE pv.profile_id = $2)`, customerID, profileID).
		Scan(&snap.Pending, &snap.Installed, &snap.Partial, &snap.Failed, &snap.Total)
	return snap, err
}

func (r *HubRepository) AssignedFolderNames(ctx context.Context, customerID, profileID int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT n.name
		FROM profile_tree_assignments a
		JOIN device_tree_nodes n ON n.id = a.tree_node_id
		WHERE a.customerid = $1 AND a.profile_id = $2
		ORDER BY lower(n.name)`, customerID, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

func (r *HubRepository) ListActivity(ctx context.Context, profileID, limit int) ([]domain.ProfileActivityEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	keys := []string{
		strconv.Itoa(profileID),
		fmt.Sprintf("profile:%d", profileID),
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, event_type, payload, created_at
		FROM domain_events
		WHERE aggregate_id = $1 OR aggregate_id = $2
		ORDER BY created_at DESC
		LIMIT $3`, keys[0], keys[1], limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ProfileActivityEvent
	for rows.Next() {
		var ev domain.ProfileActivityEvent
		var payload []byte
		if err := rows.Scan(&ev.ID, &ev.EventType, &payload, &ev.OccurredAt); err != nil {
			return nil, err
		}
		ev.SummaryKey, ev.SummaryParams, ev.ActorUserID = activitySummary(ev.EventType, payload)
		out = append(out, ev)
	}
	return out, rows.Err()
}

func activitySummary(eventType string, payload []byte) (key string, params map[string]any, actor *int) {
	params = map[string]any{}
	var raw map[string]any
	_ = json.Unmarshal(payload, &raw)
	if v, ok := raw["userId"].(float64); ok && v > 0 {
		u := int(v)
		actor = &u
	}
	switch eventType {
	case "ProfilePublished":
		key = "profile.activity.published"
		if n, ok := raw["versionNumber"].(float64); ok {
			params["versionNumber"] = int(n)
		}
	case "ProfileAssignmentChanged":
		key = "profile.activity.assigned"
		if n, ok := raw["folderName"].(string); ok {
			params["folderName"] = n
		}
		if n, ok := raw["versionNumber"].(float64); ok {
			params["versionNumber"] = int(n)
		}
	case "ProfileEnabled":
		key = "profile.activity.enabled"
	case "ProfileDisabled":
		key = "profile.activity.disabled"
	default:
		key = "profile.activity.generic"
		params["eventType"] = eventType
	}
	return key, params, actor
}

func (r *HubRepository) PinnedFromPublished(ctx context.Context, publishedVersionID int) (domain.ProfilePinned, error) {
	var pinned domain.ProfilePinned
	var settings []byte
	var publishedAt sql.NullTime
	var mainAppID sql.NullInt64
	var appCount int
	var mainAppName sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT pv.settingsjson, pv.published_at, pv.mainappid,
		       COALESCE((
		           SELECT COUNT(*)::int FROM profile_version_applications pva
		           WHERE pva.profile_version_id = pv.id AND COALESCE(pva.action, 1) = 1
		       ), 0),
		       (SELECT a.name FROM applications a
		        JOIN applicationversions av ON av.applicationid = a.id
		        WHERE av.id = pv.mainappid
		        LIMIT 1)
		FROM profile_versions pv WHERE pv.id = $1`, publishedVersionID).
		Scan(&settings, &publishedAt, &mainAppID, &appCount, &mainAppName)
	if err == sql.ErrNoRows {
		return pinned, nil
	}
	if err != nil {
		return pinned, err
	}
	if publishedAt.Valid {
		t := publishedAt.Time
		pinned.LastPublishedAt = &t
	}
	pinned.AppCount = appCount
	if mainAppName.Valid && strings.TrimSpace(mainAppName.String) != "" {
		pinned.MainAppName = strings.TrimSpace(mainAppName.String)
	}
	var doc map[string]any
	if json.Unmarshal(settings, &doc) == nil {
		if k, ok := doc["kioskMode"].(bool); ok {
			pinned.KioskMode = k
		}
		if pinned.AppCount == 0 {
			if apps, ok := doc["applications"].([]any); ok {
				pinned.AppCount = len(apps)
			}
		}
	}
	return pinned, nil
}
