package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/locations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/locations/port"
)

// Repository implements port.Repository for PostgreSQL.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

var _ port.Repository = (*Repository)(nil)

func (r *Repository) InsertBatch(ctx context.Context, records []domain.LocationRecord) error {
	if len(records) == 0 {
		return nil
	}

	const cols = "(device_id, latitude, longitude, accuracy, speed, altitude, battery_level, network_type, tracking_mode, timestamp, received_at, month)"
	values := make([]string, 0, len(records))
	args := make([]any, 0, len(records)*12)

	for i, rec := range records {
		base := i * 12
		values = append(values, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, NOW(), $%d)",
			base+1, base+2, base+3, base+4, base+5, base+6, base+7, base+8, base+9, base+10, base+11,
		))
		args = append(args,
			rec.DeviceID,
			rec.Latitude,
			rec.Longitude,
			rec.Accuracy,
			rec.Speed,
			rec.Altitude,
			rec.BatteryLevel,
			rec.NetworkType,
			rec.TrackingMode,
			rec.Timestamp,
			rec.Month,
		)
	}

	query := fmt.Sprintf(
		"INSERT INTO device_locations %s VALUES %s",
		cols, strings.Join(values, ", "),
	)
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repository) GetLastRecord(ctx context.Context, deviceID int) (*domain.LocationRecord, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, device_id, latitude, longitude, accuracy, speed, altitude, battery_level,
			network_type, tracking_mode, timestamp
		FROM device_locations
		WHERE device_id = $1
		ORDER BY timestamp DESC
		LIMIT 1`, deviceID)

	var rec domain.LocationRecord
	err := row.Scan(
		&rec.ID, &rec.DeviceID, &rec.Latitude, &rec.Longitude,
		&rec.Accuracy, &rec.Speed, &rec.Altitude, &rec.BatteryLevel,
		&rec.NetworkType, &rec.TrackingMode, &rec.Timestamp,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *Repository) QueryHistory(ctx context.Context, deviceID int, from, to int64, limit int) ([]domain.LocationRecord, int, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, device_id, latitude, longitude, accuracy, speed, altitude, battery_level,
			network_type, tracking_mode, timestamp
		FROM device_locations
		WHERE device_id = $1 AND timestamp >= $2 AND timestamp <= $3
		ORDER BY timestamp ASC
		LIMIT $4`, deviceID, from, to, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []domain.LocationRecord
	for rows.Next() {
		var rec domain.LocationRecord
		if err := rows.Scan(
			&rec.ID, &rec.DeviceID, &rec.Latitude, &rec.Longitude,
			&rec.Accuracy, &rec.Speed, &rec.Altitude, &rec.BatteryLevel,
			&rec.NetworkType, &rec.TrackingMode, &rec.Timestamp,
		); err != nil {
			return nil, 0, err
		}
		records = append(records, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count for the range.
	var total int
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM device_locations
		WHERE device_id = $1 AND timestamp >= $2 AND timestamp <= $3`,
		deviceID, from, to).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (r *Repository) QueryArchives(ctx context.Context, deviceID int, from, to int64) ([]domain.LocationArchive, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, device_id, hour_start, start_latitude, start_longitude,
			end_latitude, end_longitude, distance_traveled, point_count
		FROM device_location_archives
		WHERE device_id = $1 AND hour_start >= $2 AND hour_start <= $3
		ORDER BY hour_start ASC`, deviceID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var archives []domain.LocationArchive
	for rows.Next() {
		var a domain.LocationArchive
		if err := rows.Scan(
			&a.ID, &a.DeviceID, &a.HourStart,
			&a.StartLatitude, &a.StartLongitude,
			&a.EndLatitude, &a.EndLongitude,
			&a.DistanceTraveled, &a.PointCount,
		); err != nil {
			return nil, err
		}
		archives = append(archives, a)
	}
	return archives, rows.Err()
}

func (r *Repository) InsertArchives(ctx context.Context, archives []domain.LocationArchive) error {
	if len(archives) == 0 {
		return nil
	}

	const cols = "(device_id, hour_start, start_latitude, start_longitude, end_latitude, end_longitude, distance_traveled, point_count)"
	values := make([]string, 0, len(archives))
	args := make([]any, 0, len(archives)*8)

	for i, a := range archives {
		base := i * 8
		values = append(values, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			base+1, base+2, base+3, base+4, base+5, base+6, base+7, base+8,
		))
		args = append(args,
			a.DeviceID, a.HourStart,
			a.StartLatitude, a.StartLongitude,
			a.EndLatitude, a.EndLongitude,
			a.DistanceTraveled, a.PointCount,
		)
	}

	query := fmt.Sprintf(
		"INSERT INTO device_location_archives %s VALUES %s ON CONFLICT (device_id, hour_start) DO NOTHING",
		cols, strings.Join(values, ", "),
	)
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repository) DeleteByIDs(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query := fmt.Sprintf("DELETE FROM device_locations WHERE id IN (%s)", strings.Join(placeholders, ","))
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repository) GetRecordsOlderThan(ctx context.Context, cutoffMs int64, batchSize int) ([]domain.LocationRecord, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, device_id, latitude, longitude, accuracy, speed, altitude, battery_level,
			network_type, tracking_mode, timestamp
		FROM device_locations
		WHERE timestamp < $1
		ORDER BY timestamp ASC
		LIMIT $2`, cutoffMs, batchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []domain.LocationRecord
	for rows.Next() {
		var rec domain.LocationRecord
		if err := rows.Scan(
			&rec.ID, &rec.DeviceID, &rec.Latitude, &rec.Longitude,
			&rec.Accuracy, &rec.Speed, &rec.Altitude, &rec.BatteryLevel,
			&rec.NetworkType, &rec.TrackingMode, &rec.Timestamp,
		); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}


