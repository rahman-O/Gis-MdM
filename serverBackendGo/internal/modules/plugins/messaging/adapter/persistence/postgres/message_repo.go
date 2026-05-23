package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/messaging/domain"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Search(ctx context.Context, customerID int64, f domain.MessageFilter) ([]domain.Message, int64, error) {
	if f.PageSize <= 0 {
		f.PageSize = 50
	}
	if f.PageNum <= 0 {
		f.PageNum = 1
	}
	var total int64
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM plugin_messaging_messages WHERE customerid = $1`, customerID).Scan(&total)
	offset := (f.PageNum - 1) * f.PageSize
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, customerid, deviceid, ts, COALESCE(message,''), status
		FROM plugin_messaging_messages WHERE customerid = $1
		ORDER BY ts DESC LIMIT $2 OFFSET $3`, customerID, f.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []domain.Message
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(&m.ID, &m.CustomerID, &m.DeviceID, &m.Ts, &m.Message, &m.Status); err != nil {
			return nil, 0, err
		}
		items = append(items, m)
	}
	return items, total, rows.Err()
}

func (r *MessageRepository) Insert(ctx context.Context, customerID, deviceID int64, text string) (int64, error) {
	ts := time.Now().UnixMilli()
	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO plugin_messaging_messages (customerid, deviceid, ts, message, status)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		customerID, deviceID, ts, text, domain.StatusPending).Scan(&id)
	return id, err
}

func (r *MessageRepository) Delete(ctx context.Context, customerID, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM plugin_messaging_messages WHERE id = $1 AND customerid = $2`, id, customerID)
	return err
}

func (r *MessageRepository) Purge(ctx context.Context, customerID int64, days int) error {
	cutoff := time.Now().AddDate(0, 0, -days).UnixMilli()
	_, err := r.db.ExecContext(ctx, `DELETE FROM plugin_messaging_messages WHERE customerid = $1 AND ts < $2`, customerID, cutoff)
	return err
}

func (r *MessageRepository) UpdateStatus(ctx context.Context, id int64, status int) error {
	res, err := r.db.ExecContext(ctx, `UPDATE plugin_messaging_messages SET status = $1 WHERE id = $2`, status, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *MessageRepository) DeviceIDByNumber(ctx context.Context, customerID int64, number string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM devices WHERE customerid = $1 AND lower(number) = lower($2)`, customerID, number).Scan(&id)
	return id, err
}

// QueuePayload builds agent textMessage payload per Java MessagingResource.
func QueuePayload(id int64, text string) string {
	escaped := ""
	for _, r := range text {
		if r == '"' || r == '\\' {
			escaped += `\`
		}
		escaped += string(r)
	}
	return fmt.Sprintf(`{id:%d,text:"%s"}`, id, escaped)
}
