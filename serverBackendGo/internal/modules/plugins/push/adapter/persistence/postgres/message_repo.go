package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/domain"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Search(ctx context.Context, customerID int64, f domain.PushMessageFilter) ([]domain.PluginPushMessage, int64, error) {
	if f.PageSize <= 0 {
		f.PageSize = 50
	}
	if f.PageNum <= 0 {
		f.PageNum = 1
	}
	offset := (f.PageNum - 1) * f.PageSize
	var total int64
	_ = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM plugin_push_messages WHERE customerid = $1`, customerID).Scan(&total)
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, customerid, deviceid, ts, COALESCE(messagetype,''), COALESCE(payload,'')
		FROM plugin_push_messages WHERE customerid = $1
		ORDER BY ts DESC LIMIT $2 OFFSET $3`, customerID, f.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []domain.PluginPushMessage
	for rows.Next() {
		var m domain.PluginPushMessage
		if err := rows.Scan(&m.ID, &m.CustomerID, &m.DeviceID, &m.Ts, &m.MessageType, &m.Payload); err != nil {
			return nil, 0, err
		}
		items = append(items, m)
	}
	return items, total, rows.Err()
}

func (r *MessageRepository) InsertHistory(ctx context.Context, customerID, deviceID int64, messageType, payload string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO plugin_push_messages (customerid, deviceid, ts, messagetype, payload)
		VALUES ($1, $2, $3, $4, $5)`, customerID, deviceID, time.Now().UnixMilli(), messageType, payload)
	return err
}

func (r *MessageRepository) Delete(ctx context.Context, customerID, id int64) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM plugin_push_messages WHERE id = $1 AND customerid = $2`, id, customerID)
	return err
}

func (r *MessageRepository) Purge(ctx context.Context, customerID int64, days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days).UnixMilli()
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM plugin_push_messages WHERE customerid = $1 AND ts < $2`, customerID, cutoff)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *MessageRepository) DeviceIDByNumber(ctx context.Context, customerID int64, number string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM devices WHERE customerid = $1 AND lower(number) = lower($2)`, customerID, number).Scan(&id)
	return id, err
}
