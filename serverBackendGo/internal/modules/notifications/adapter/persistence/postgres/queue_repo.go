package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/notifications/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/notifications/port"
)

type QueueRepository struct {
	db *sql.DB
}

func NewQueueRepository(db *sql.DB) *QueueRepository {
	return &QueueRepository{db: db}
}

var _ port.MessageQueue = (*QueueRepository)(nil)

func (r *QueueRepository) Enqueue(ctx context.Context, deviceID int64, messageType, payload string) error {
	var msgID int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO pushmessages (messagetype, deviceid, payload)
		VALUES ($1, $2, $3)
		RETURNING id`, messageType, deviceID, payload).Scan(&msgID)
	if err != nil {
		return err
	}
	now := time.Now().UnixMilli()
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO pendingpushes (messageid, status, createtime)
		VALUES ($1, 0, $2)`, msgID, now)
	return err
}

func (r *QueueRepository) ListPendingForDevice(ctx context.Context, deviceID int64) ([]domain.PlainPushMessage, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT pm.id, pm.messagetype, COALESCE(pm.payload, '')
		FROM pushmessages pm
		JOIN pendingpushes pp ON pp.messageid = pm.id
		WHERE pm.deviceid = $1 AND pp.status = 0
		ORDER BY pm.id`, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.PlainPushMessage
	for rows.Next() {
		var m domain.PlainPushMessage
		if err := rows.Scan(&m.ID, &m.MessageType, &m.Payload); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *QueueRepository) MarkDelivered(ctx context.Context, messageIDs []int64) error {
	if len(messageIDs) == 0 {
		return nil
	}
	now := time.Now().UnixMilli()
	for _, id := range messageIDs {
		_, err := r.db.ExecContext(ctx, `
			UPDATE pendingpushes SET status = 1, sendtime = $1
			WHERE messageid = $2 AND status = 0`, now, id)
		if err != nil {
			return err
		}
	}
	return nil
}

type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

var _ port.DeviceLookup = (*DeviceRepository)(nil)

func (r *DeviceRepository) DeviceIDByNumber(ctx context.Context, deviceNumber string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM devices WHERE lower(number) = lower($1)
		UNION ALL
		SELECT id FROM devices WHERE oldnumber IS NOT NULL AND lower(oldnumber) = lower($1)
		LIMIT 1`, deviceNumber).Scan(&id)
	return id, err
}
