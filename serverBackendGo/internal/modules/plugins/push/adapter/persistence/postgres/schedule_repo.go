package postgres

import (
	"context"
	"database/sql"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/domain"
)

type ScheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepository(db *sql.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Search(ctx context.Context, customerID int64, f domain.PushScheduleFilter) ([]domain.PluginPushSchedule, int64, error) {
	if f.PageSize <= 0 {
		f.PageSize = 50
	}
	if f.PageNum <= 0 {
		f.PageNum = 1
	}
	var total int64
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM plugin_push_schedule WHERE customerid = $1`, customerID).Scan(&total)
	offset := (f.PageNum - 1) * f.PageSize
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, customerid, deviceid, groupid, configurationid, COALESCE(scope,''),
			COALESCE(messagetype,''), COALESCE(payload,''), COALESCE(comment,''),
			COALESCE(min,''), COALESCE(hour,''), COALESCE(day,''), COALESCE(weekday,''), COALESCE(month,'')
		FROM plugin_push_schedule WHERE customerid = $1 ORDER BY id DESC LIMIT $2 OFFSET $3`,
		customerID, f.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []domain.PluginPushSchedule
	for rows.Next() {
		var s domain.PluginPushSchedule
		if err := rows.Scan(&s.ID, &s.CustomerID, &s.DeviceID, &s.GroupID, &s.ConfigurationID, &s.Scope,
			&s.MessageType, &s.Payload, &s.Comment, &s.Min, &s.Hour, &s.Day, &s.Weekday, &s.Month); err != nil {
			return nil, 0, err
		}
		items = append(items, s)
	}
	return items, total, rows.Err()
}

func (r *ScheduleRepository) Save(ctx context.Context, s domain.PluginPushSchedule) (int64, error) {
	if s.ID == 0 {
		var id int64
		err := r.db.QueryRowContext(ctx, `
			INSERT INTO plugin_push_schedule (customerid, deviceid, groupid, configurationid, scope, messagetype, payload, comment, min, hour, day, weekday, month)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`,
			s.CustomerID, s.DeviceID, s.GroupID, s.ConfigurationID, s.Scope, s.MessageType, s.Payload, s.Comment,
			s.Min, s.Hour, s.Day, s.Weekday, s.Month).Scan(&id)
		return id, err
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE plugin_push_schedule SET deviceid=$2, groupid=$3, configurationid=$4, scope=$5, messagetype=$6,
			payload=$7, comment=$8, min=$9, hour=$10, day=$11, weekday=$12, month=$13
		WHERE id=$1 AND customerid=$14`,
		s.ID, s.DeviceID, s.GroupID, s.ConfigurationID, s.Scope, s.MessageType, s.Payload, s.Comment,
		s.Min, s.Hour, s.Day, s.Weekday, s.Month, s.CustomerID)
	return s.ID, err
}

func (r *ScheduleRepository) Delete(ctx context.Context, customerID, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM plugin_push_schedule WHERE id = $1 AND customerid = $2`, id, customerID)
	return err
}

func (r *ScheduleRepository) ResolveDeviceID(ctx context.Context, customerID int64, number string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM devices WHERE customerid = $1 AND lower(number) = lower($2)`, customerID, number).Scan(&id)
	return id, err
}

// ListAll returns every schedule row (filtered in application by cron masks).
func (r *ScheduleRepository) ListAll(ctx context.Context) ([]domain.PluginPushSchedule, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, customerid, deviceid, groupid, configurationid, COALESCE(scope,''),
			COALESCE(messagetype,''), COALESCE(payload,''), COALESCE(comment,''),
			COALESCE(min,'*'), COALESCE(hour,'*'), COALESCE(day,'*'), COALESCE(weekday,'*'), COALESCE(month,'*')
		FROM plugin_push_schedule`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.PluginPushSchedule
	for rows.Next() {
		var s domain.PluginPushSchedule
		if err := rows.Scan(&s.ID, &s.CustomerID, &s.DeviceID, &s.GroupID, &s.ConfigurationID, &s.Scope,
			&s.MessageType, &s.Payload, &s.Comment, &s.Min, &s.Hour, &s.Day, &s.Weekday, &s.Month); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, rows.Err()
}

// ResolveDeviceIDs returns target device IDs for a scheduled task scope.
func (r *ScheduleRepository) ResolveDeviceIDs(ctx context.Context, task domain.PluginPushSchedule) ([]int64, error) {
	switch task.Scope {
	case "device":
		if task.DeviceID > 0 {
			return []int64{task.DeviceID}, nil
		}
		return nil, nil
	case "group":
		return r.deviceIDsByFilter(ctx, `SELECT id FROM devices WHERE customerid = $1 AND groupid = $2`,
			task.CustomerID, task.GroupID)
	case "configuration":
		return r.deviceIDsByFilter(ctx, `SELECT id FROM devices WHERE customerid = $1 AND configurationid = $2`,
			task.CustomerID, task.ConfigurationID)
	default:
		return r.deviceIDsByFilter(ctx, `SELECT id FROM devices WHERE customerid = $1`, task.CustomerID)
	}
}

func (r *ScheduleRepository) deviceIDsByFilter(ctx context.Context, q string, args ...any) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
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
