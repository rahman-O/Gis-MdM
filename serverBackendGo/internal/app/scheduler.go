package app

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/config"
	notifpostgres "github.com/gis-mdm/server-backend-go/internal/modules/notifications/adapter/persistence/postgres"
	pluginpostgres "github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/adapter/persistence/postgres"
	pluginapp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/application"
)

// StartPushScheduleRunner runs plugin push cron until ctx is cancelled.
func StartPushScheduleRunner(ctx context.Context, cfg config.Config, db *sql.DB, log *slog.Logger) {
	if db == nil || !cfg.ModulePluginsEnabled || !cfg.IsPluginEnabled("push") {
		return
	}
	interval := time.Duration(cfg.PushScheduleIntervalSec) * time.Second
	if interval < time.Second {
		interval = time.Minute
	}
	runner := pluginapp.NewScheduleRunner(
		pluginpostgres.NewScheduleRepository(db),
		notifpostgres.NewQueueRepository(db),
		pluginpostgres.NewMessageRepository(db),
		log,
	)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runner.RunOnce(ctx)
			}
		}
	}()
	if log != nil {
		log.Info("push schedule runner started", slog.Duration("interval", interval))
	}
}
