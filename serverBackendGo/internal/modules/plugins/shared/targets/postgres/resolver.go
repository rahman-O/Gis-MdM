package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/shared/targets"
)

type Resolver struct {
	db *sql.DB
}

func NewResolver(db *sql.DB) *Resolver {
	return &Resolver{db: db}
}

var _ targets.Resolver = (*Resolver)(nil)

func (r *Resolver) userAllDevices(ctx context.Context, userID int64) (bool, error) {
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

func (r *Resolver) deviceIDs(ctx context.Context, customerID, userID int64, extra string, args []any) ([]int64, error) {
	allDevices, err := r.userAllDevices(ctx, userID)
	if err != nil {
		return nil, err
	}
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

func (r *Resolver) DeviceIDsByNumbers(ctx context.Context, customerID, userID int64, numbers []string) ([]int64, error) {
	if len(numbers) == 0 {
		return nil, nil
	}
	lower := make([]string, len(numbers))
	for i, n := range numbers {
		lower[i] = n
	}
	return r.deviceIDs(ctx, customerID, userID, `AND lower(d.number) = ANY($4)`, []any{pq.Array(lower)})
}

func (r *Resolver) DeviceIDsByGroupID(ctx context.Context, customerID, userID, groupID int64) ([]int64, error) {
	return r.deviceIDs(ctx, customerID, userID, `AND EXISTS (
		SELECT 1 FROM devicegroups dg WHERE dg.deviceid = d.id AND dg.groupid = $4)`, []any{groupID})
}

func (r *Resolver) DeviceIDsByConfigurationID(ctx context.Context, customerID, userID, configurationID int64) ([]int64, error) {
	return r.deviceIDs(ctx, customerID, userID, `AND d.configurationid = $4`, []any{configurationID})
}

func (r *Resolver) AllDeviceIDs(ctx context.Context, customerID, userID int64) ([]int64, error) {
	return r.deviceIDs(ctx, customerID, userID, "", nil)
}
