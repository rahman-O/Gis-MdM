package application

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"time"

	pushport "github.com/gis-mdm/server-backend-go/internal/platform/push/port"
)

// DomainEventsWorker processes domain_events outbox rows (debounced poll).
type DomainEventsWorker struct {
	db       *sql.DB
	queue    pushport.MessageQueue
	interval time.Duration
	log      *slog.Logger
}

func NewDomainEventsWorker(db *sql.DB, queue pushport.MessageQueue, log *slog.Logger) *DomainEventsWorker {
	if log == nil {
		log = slog.Default()
	}
	return &DomainEventsWorker{db: db, queue: queue, interval: 5 * time.Second, log: log}
}

func (w *DomainEventsWorker) Start(ctx context.Context) {
	if w == nil || w.db == nil || w.queue == nil {
		return
	}
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.processBatch(ctx)
			}
		}
	}()
}

func (w *DomainEventsWorker) processBatch(ctx context.Context) {
	rows, err := w.db.QueryContext(ctx, `
		SELECT id, event_type, payload FROM domain_events
		WHERE processed_at IS NULL
		ORDER BY id ASC LIMIT 20`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var eventType string
		var payload []byte
		if err := rows.Scan(&id, &eventType, &payload); err != nil {
			continue
		}
		switch eventType {
		case "ProfilePublished":
			w.handleProfilePublished(ctx, payload)
		}
		_, _ = w.db.ExecContext(ctx, `UPDATE domain_events SET processed_at = NOW() WHERE id = $1`, id)
	}
}

func (w *DomainEventsWorker) handleProfilePublished(ctx context.Context, payload []byte) {
	var body struct {
		ProfileID int `json:"profileId"`
	}
	if json.Unmarshal(payload, &body) != nil || body.ProfileID <= 0 {
		return
	}
	rows, err := w.db.QueryContext(ctx, `
		SELECT DISTINCT d.id FROM devices d
		JOIN enrollment_routes er ON er.id = d.enrollment_route_id
		JOIN profile_versions pv ON pv.id = er.profile_version_id
		WHERE pv.profile_id = $1`, body.ProfileID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var deviceID int64
		if err := rows.Scan(&deviceID); err != nil {
			continue
		}
		_ = w.queue.Enqueue(ctx, deviceID, TypeConfigUpdated, "")
	}
}
