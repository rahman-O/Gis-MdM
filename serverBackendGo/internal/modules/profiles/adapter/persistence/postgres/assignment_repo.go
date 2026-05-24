package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/port"
)

// AssignmentRepository implements port.RolloutStore.
type AssignmentRepository struct {
	db *sql.DB
}

func NewAssignmentRepository(db *sql.DB) *AssignmentRepository {
	return &AssignmentRepository{db: db}
}

var _ port.RolloutStore = (*AssignmentRepository)(nil)

func (r *AssignmentRepository) CountSubtreeDevices(ctx context.Context, customerID, treeNodeID int) (int, error) {
	var path string
	err := r.db.QueryRowContext(ctx, `
		SELECT path FROM device_tree_nodes WHERE id = $1 AND customerid = $2`, treeNodeID, customerID).Scan(&path)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	var n int
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)::int FROM devices d
		JOIN device_tree_nodes n ON n.id = d.tree_node_id
		WHERE d.customerid = $1 AND n.path LIKE $2 || '%'`, customerID, path).Scan(&n)
	return n, err
}

func (r *AssignmentRepository) ListAssignments(ctx context.Context, customerID, profileID int) ([]domain.ProfileTreeAssignment, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT a.id, a.tree_node_id, n.name, n.path, a.profile_version_id, pv.version_number,
		       (SELECT COUNT(*)::int FROM devices d
		        JOIN device_tree_nodes dn ON dn.id = d.tree_node_id
		        WHERE d.customerid = $1 AND dn.path LIKE n.path || '%'),
		       a.created_at
		FROM profile_tree_assignments a
		JOIN device_tree_nodes n ON n.id = a.tree_node_id
		JOIN profile_versions pv ON pv.id = a.profile_version_id
		WHERE a.customerid = $1 AND a.profile_id = $2
		ORDER BY n.path`, customerID, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ProfileTreeAssignment
	for rows.Next() {
		var item domain.ProfileTreeAssignment
		if err := rows.Scan(&item.AssignmentID, &item.TreeNodeID, &item.TreeNodeName, &item.TreePath,
			&item.ProfileVersionID, &item.VersionNumber, &item.DeviceCount, &item.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *AssignmentRepository) GetAssignmentImpact(ctx context.Context, customerID, treeNodeID int) (int, string, error) {
	var name string
	err := r.db.QueryRowContext(ctx, `
		SELECT name FROM device_tree_nodes WHERE id = $1 AND customerid = $2`, treeNodeID, customerID).Scan(&name)
	if err == sql.ErrNoRows {
		return 0, "", sql.ErrNoRows
	}
	if err != nil {
		return 0, "", err
	}
	n, err := r.CountSubtreeDevices(ctx, customerID, treeNodeID)
	return n, name, err
}

func (r *AssignmentRepository) UpsertAssignment(ctx context.Context, customerID, profileID, profileVersionID, treeNodeID, createdBy int) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO profile_tree_assignments (customerid, profile_id, profile_version_id, tree_node_id, created_by)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (customerid, tree_node_id) DO UPDATE SET
			profile_id = EXCLUDED.profile_id,
			profile_version_id = EXCLUDED.profile_version_id,
			created_at = now(),
			created_by = EXCLUDED.created_by
		RETURNING id`, customerID, profileID, profileVersionID, treeNodeID, nullIntSQL(createdBy)).Scan(&id)
	return id, err
}

func (r *AssignmentRepository) DeleteAssignment(ctx context.Context, customerID, profileID, assignmentID int) error {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM profile_tree_assignments
		WHERE id = $1 AND customerid = $2 AND profile_id = $3`, assignmentID, customerID, profileID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *AssignmentRepository) MarkSubtreePending(ctx context.Context, customerID, treeNodeID, targetVersionID int) (int, error) {
	return r.markSubtreePending(ctx, r.db, customerID, treeNodeID, targetVersionID)
}

func (r *AssignmentRepository) markSubtreePending(ctx context.Context, exec interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}, customerID, treeNodeID, targetVersionID int) (int, error) {
	var path string
	if err := exec.QueryRowContext(ctx, `
		SELECT path FROM device_tree_nodes WHERE id = $1 AND customerid = $2`, treeNodeID, customerID).Scan(&path); err != nil {
		return 0, err
	}
	res, err := exec.ExecContext(ctx, `
		UPDATE devices d SET
			target_profile_version_id = $3,
			profile_rollout_status = 'pending',
			profile_rollout_reason = NULL,
			profile_rollout_updated_at = now()
		FROM device_tree_nodes n
		WHERE d.tree_node_id = n.id AND d.customerid = $1 AND n.path LIKE $2 || '%'`,
		customerID, path, targetVersionID)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}

func (r *AssignmentRepository) ListAssignmentsForPublishImpact(ctx context.Context, customerID, profileID int) ([]domain.PublishImpactAssignment, error) {
	list, err := r.ListAssignments(ctx, customerID, profileID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.PublishImpactAssignment, 0, len(list))
	for _, a := range list {
		out = append(out, domain.PublishImpactAssignment{
			AssignmentID:         a.AssignmentID,
			TreeNodeID:           a.TreeNodeID,
			TreeNodeName:         a.TreeNodeName,
			CurrentVersionNumber: a.VersionNumber,
			DeviceCount:          a.DeviceCount,
		})
	}
	return out, nil
}

func (r *AssignmentRepository) BumpAllAssignmentsOnPublish(ctx context.Context, customerID, profileID, newVersionID int) (int, int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()
	updated, affected, err := r.bumpAllAssignmentsOnPublishTx(ctx, tx, customerID, profileID, newVersionID)
	if err != nil {
		return 0, 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}
	return updated, affected, nil
}

// BumpAllAssignmentsOnPublishTx bumps folder assignments inside an existing transaction.
func (r *AssignmentRepository) BumpAllAssignmentsOnPublishTx(ctx context.Context, tx *sql.Tx, customerID, profileID, newVersionID int) (int, int, error) {
	return r.bumpAllAssignmentsOnPublishTx(ctx, tx, customerID, profileID, newVersionID)
}

func (r *AssignmentRepository) bumpAllAssignmentsOnPublishTx(ctx context.Context, tx *sql.Tx, customerID, profileID, newVersionID int) (int, int, error) {
	res, err := tx.ExecContext(ctx, `
		UPDATE profile_tree_assignments
		SET profile_version_id = $3
		WHERE customerid = $1 AND profile_id = $2 AND profile_version_id IS DISTINCT FROM $3`,
		customerID, profileID, newVersionID)
	if err != nil {
		return 0, 0, err
	}
	updated, _ := res.RowsAffected()
	rows, err := tx.QueryContext(ctx, `
		SELECT tree_node_id FROM profile_tree_assignments
		WHERE customerid = $1 AND profile_id = $2`, customerID, profileID)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()
	var nodes []int
	for rows.Next() {
		var nodeID int
		if err := rows.Scan(&nodeID); err != nil {
			return 0, 0, err
		}
		nodes = append(nodes, nodeID)
	}
	if err := rows.Err(); err != nil {
		return 0, 0, err
	}
	devicesAffected := 0
	for _, nodeID := range nodes {
		n, err := r.markSubtreePending(ctx, tx, customerID, nodeID, newVersionID)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return 0, 0, err
		}
		devicesAffected += n
	}
	return int(updated), devicesAffected, nil
}

func (r *AssignmentRepository) ListRolloutDevices(ctx context.Context, customerID, profileID int, q domain.RolloutDevicesQuery) ([]domain.DeviceRolloutRow, int, error) {
	if q.PageSize <= 0 {
		q.PageSize = 25
	}
	if q.Page < 1 {
		q.Page = 1
	}
	where := `WHERE d.customerid = $1 AND (
		d.target_profile_version_id IN (SELECT profile_version_id FROM profile_tree_assignments WHERE profile_id = $2 AND customerid = $1)
		OR d.applied_profile_version_id IN (SELECT id FROM profile_versions WHERE profile_id = $2)
	)`
	args := []any{customerID, profileID}
	argN := 3
	if q.TreeNodeID != nil && *q.TreeNodeID > 0 {
		var path string
		if err := r.db.QueryRowContext(ctx, `
			SELECT path FROM device_tree_nodes WHERE id = $1 AND customerid = $2`, *q.TreeNodeID, customerID).Scan(&path); err == nil {
			where += fmt.Sprintf(` AND EXISTS (
				SELECT 1 FROM device_tree_nodes tn WHERE tn.id = d.tree_node_id AND tn.path LIKE $%d || '%%')`, argN)
			args = append(args, path)
			argN++
		}
	}
	if st := strings.TrimSpace(q.Status); st != "" {
		where += fmt.Sprintf(` AND d.profile_rollout_status = $%d`, argN)
		args = append(args, st)
		argN++
	}
	var total int
	countQ := `SELECT COUNT(*)::int FROM devices d ` + where
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (q.Page - 1) * q.PageSize
	listQ := `
		SELECT d.id, COALESCE(NULLIF(TRIM(d.description), ''), d.number),
		       d.tree_node_id, COALESCE(tn.name, ''),
		       d.target_profile_version_id, tv.version_number,
		       d.applied_profile_version_id, av.version_number,
		       COALESCE(d.profile_rollout_status, 'pending'),
		       COALESCE(d.profile_rollout_reason, ''),
		       d.lastupdate
		FROM devices d
		LEFT JOIN device_tree_nodes tn ON tn.id = d.tree_node_id
		LEFT JOIN profile_versions tv ON tv.id = d.target_profile_version_id
		LEFT JOIN profile_versions av ON av.id = d.applied_profile_version_id
		` + where + fmt.Sprintf(` ORDER BY d.lastupdate DESC NULLS LAST LIMIT $%d OFFSET $%d`, argN, argN+1)
	args = append(args, q.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, listQ, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []domain.DeviceRolloutRow
	for rows.Next() {
		var row domain.DeviceRolloutRow
		var treeID sql.NullInt64
		var tgtID, appID sql.NullInt64
		var tgtNum, appNum sql.NullInt64
		var last sql.NullInt64
		if err := rows.Scan(&row.DeviceID, &row.DeviceName, &treeID, &row.TreeNodeName,
			&tgtID, &tgtNum, &appID, &appNum, &row.Status, &row.Reason, &last); err != nil {
			return nil, 0, err
		}
		if treeID.Valid {
			v := int(treeID.Int64)
			row.TreeNodeID = &v
		}
		if tgtID.Valid {
			v := int(tgtID.Int64)
			row.TargetVersionID = &v
		}
		if tgtNum.Valid {
			v := int(tgtNum.Int64)
			row.TargetVersionNumber = &v
		}
		if appID.Valid {
			v := int(appID.Int64)
			row.AppliedVersionID = &v
		}
		if appNum.Valid {
			v := int(appNum.Int64)
			row.AppliedVersionNumber = &v
		}
		if last.Valid {
			v := last.Int64
			row.LastUpdate = &v
		}
		items = append(items, row)
	}
	return items, total, rows.Err()
}

func (r *AssignmentRepository) IsVersionPublished(ctx context.Context, customerID, profileID, versionID int) (bool, error) {
	var ok bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM profile_versions pv
			JOIN profiles p ON p.id = pv.profile_id
			WHERE pv.id = $1 AND p.id = $2 AND p.customerid = $3 AND pv.status = 'published'
		)`, versionID, profileID, customerID).Scan(&ok)
	return ok, err
}

func (r *AssignmentRepository) IsProfileEnabled(ctx context.Context, customerID, profileID int) (bool, error) {
	var enabled bool
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(enabled, true) FROM profiles WHERE id = $1 AND customerid = $2`, profileID, customerID).Scan(&enabled)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return enabled, err
}

func (r *AssignmentRepository) SetProfileEnabled(ctx context.Context, customerID, profileID int, enabled bool) (int, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE profiles SET enabled = $3 WHERE id = $1 AND customerid = $2`, profileID, customerID, enabled)
	if err != nil {
		return 0, err
	}
	if !enabled {
		return 0, nil
	}
	return r.MarkProfileDevicesPending(ctx, customerID, profileID)
}

func (r *AssignmentRepository) MarkProfileDevicesPending(ctx context.Context, customerID, profileID int) (int, error) {
	res, err := r.db.ExecContext(ctx, `
		UPDATE devices d SET
			profile_rollout_status = 'pending',
			profile_rollout_reason = 'Profile re-enabled',
			profile_rollout_updated_at = now()
		WHERE d.customerid = $1 AND (
			d.target_profile_version_id IN (SELECT profile_version_id FROM profile_tree_assignments WHERE profile_id = $2)
			OR d.enrollment_route_id IN (
				SELECT er.id FROM enrollment_routes er
				JOIN profile_versions pv ON pv.id = er.profile_version_id
				WHERE pv.profile_id = $2
			)
		)`, customerID, profileID)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}

func nullIntSQL(v int) any {
	if v <= 0 {
		return nil
	}
	return v
}

// DeviceContext for resolver.
type DeviceContext struct {
	DeviceID          int64
	CustomerID        int
	TreeNodeID        sql.NullInt64
	EnrollmentRouteID sql.NullInt64
}

func (r *AssignmentRepository) LoadDeviceContext(ctx context.Context, deviceID int64) (*DeviceContext, error) {
	var dc DeviceContext
	err := r.db.QueryRowContext(ctx, `
		SELECT id, customerid, tree_node_id, enrollment_route_id
		FROM devices WHERE id = $1`, deviceID).Scan(
		&dc.DeviceID, &dc.CustomerID, &dc.TreeNodeID, &dc.EnrollmentRouteID)
	if err != nil {
		return nil, err
	}
	return &dc, nil
}

func (r *AssignmentRepository) ResolveTreeVersion(ctx context.Context, customerID int, treeNodeID int64) (profileID, versionID int, ok bool, err error) {
	if treeNodeID <= 0 {
		return 0, 0, false, nil
	}
	var path string
	if err := r.db.QueryRowContext(ctx, `
		SELECT path FROM device_tree_nodes WHERE id = $1 AND customerid = $2`, treeNodeID, customerID).Scan(&path); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, false, nil
		}
		return 0, 0, false, err
	}
	// Walk ancestors deepest first via path prefix match on assignments
	err = r.db.QueryRowContext(ctx, `
		SELECT a.profile_id, a.profile_version_id
		FROM profile_tree_assignments a
		JOIN device_tree_nodes n ON n.id = a.tree_node_id
		JOIN profiles p ON p.id = a.profile_id
		WHERE a.customerid = $1 AND p.enabled = true
		  AND $2 LIKE n.path || '%'
		ORDER BY n.depth DESC
		LIMIT 1`, customerID, path).Scan(&profileID, &versionID)
	if err == sql.ErrNoRows {
		return 0, 0, false, nil
	}
	if err != nil {
		return 0, 0, false, err
	}
	return profileID, versionID, true, nil
}

func (r *AssignmentRepository) SetDeviceTargetVersion(ctx context.Context, deviceID int64, versionID int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE devices SET target_profile_version_id = $2,
			profile_rollout_status = 'pending',
			profile_rollout_updated_at = now()
		WHERE id = $1`, deviceID, nullIntSQL(versionID))
	return err
}

func (r *AssignmentRepository) SetDeviceAppliedVersion(ctx context.Context, deviceID int64, versionID int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE devices SET applied_profile_version_id = $2, profile_rollout_updated_at = now()
		WHERE id = $1`, deviceID, versionID)
	return err
}

func (r *AssignmentRepository) UpdateRolloutFields(ctx context.Context, deviceID int64, status, reason string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE devices SET profile_rollout_status = $2, profile_rollout_reason = $3, profile_rollout_updated_at = $4
		WHERE id = $1`, deviceID, status, reason, time.Now())
	return err
}
