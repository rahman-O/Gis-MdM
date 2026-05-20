package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/gis-mdm/server-backend-go/internal/modules/push/port"
)

type TargetRepository struct {
	db *sql.DB
}

func NewTargetRepository(db *sql.DB) *TargetRepository {
	return &TargetRepository{db: db}
}

var _ port.TargetResolver = (*TargetRepository)(nil)

func (r *TargetRepository) userAllDevices(ctx context.Context, userID int64) (bool, error) {
	var all bool
	err := r.db.QueryRowContext(ctx, `SELECT alldevicesavailable FROM users WHERE id = $1`, userID).Scan(&all)
	return all, err
}

const accessWhere = `
	d.customerid = $1
	AND ($2 = TRUE OR EXISTS (
		SELECT 1 FROM devicegroups dg
		JOIN groups g ON g.id = dg.groupid
		LEFT JOIN userdevicegroupsaccess a ON a.groupid = g.id AND a.userid = $3
		WHERE dg.deviceid = d.id AND ($2 = TRUE OR a.groupid IS NOT NULL)
	))`

func (r *TargetRepository) deviceIDs(ctx context.Context, customerID, userID int64, allDevices bool, extra string, args []any) ([]int64, error) {
	base := []any{customerID, allDevices, userID}
	base = append(base, args...)
	q := fmt.Sprintf(`SELECT DISTINCT d.id FROM devices d WHERE %s %s`, accessWhere, extra)
	rows, err := r.db.QueryContext(ctx, q, base...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *TargetRepository) DeviceIDsByNumbers(ctx context.Context, customerID, userID int64, _ bool, numbers []string) ([]int64, error) {
	allDevices, err := r.userAllDevices(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(numbers) == 0 {
		return nil, nil
	}
	lower := make([]string, len(numbers))
	for i, n := range numbers {
		lower[i] = strings.ToLower(strings.TrimSpace(n))
	}
	return r.deviceIDs(ctx, customerID, userID, allDevices, `AND lower(d.number) = ANY($4)`, []any{pq.Array(lower)})
}

func (r *TargetRepository) DeviceIDsByGroupNames(ctx context.Context, customerID, userID int64, _ bool, groups []string) ([]int64, error) {
	allDevices, err := r.userAllDevices(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return nil, nil
	}
	names := make([]string, len(groups))
	for i, g := range groups {
		names[i] = strings.ToLower(strings.TrimSpace(g))
	}
	return r.deviceIDs(ctx, customerID, userID, allDevices, `
		AND EXISTS (
			SELECT 1 FROM devicegroups dg2
			JOIN groups g2 ON g2.id = dg2.groupid
			WHERE dg2.deviceid = d.id AND lower(g2.name) = ANY($4)
		)`, []any{pq.Array(names)})
}

func (r *TargetRepository) AllDeviceIDs(ctx context.Context, customerID, userID int64, _ bool) ([]int64, error) {
	allDevices, err := r.userAllDevices(ctx, userID)
	if err != nil {
		return nil, err
	}
	return r.deviceIDs(ctx, customerID, userID, allDevices, "", nil)
}
