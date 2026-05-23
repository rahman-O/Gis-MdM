package postgres

import (
	"context"
	"database/sql"
)

// DeviceLookup resolves device IDs for push notifications.
type DeviceLookup struct {
	db *sql.DB
}

func NewDeviceLookup(db *sql.DB) *DeviceLookup {
	return &DeviceLookup{db: db}
}

func (r *DeviceLookup) DeviceIDsByConfiguration(ctx context.Context, configurationID int) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id FROM devices WHERE configurationid = $1`, configurationID)
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
