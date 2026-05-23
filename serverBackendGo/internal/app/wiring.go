package app

import (
	"database/sql"
	"log/slog"

	"github.com/gis-mdm/server-backend-go/internal/config"
	"github.com/gis-mdm/server-backend-go/internal/module"
	notifpostgres "github.com/gis-mdm/server-backend-go/internal/modules/notifications/adapter/persistence/postgres"
	pushapp "github.com/gis-mdm/server-backend-go/internal/platform/push/application"
	pushpostgres "github.com/gis-mdm/server-backend-go/internal/platform/push/adapter/postgres"
)

func moduleDependencies(cfg config.Config, db *sql.DB, log *slog.Logger) module.Dependencies {
	deps := module.Dependencies{Config: cfg, DB: db, Log: log}
	if !cfg.PushNotifierEnabled || db == nil {
		return deps
	}
	queue := notifpostgres.NewQueueRepository(db)
	lookup := pushpostgres.NewDeviceLookup(db)
	n := pushapp.NewNotifier(queue, lookup, log)
	deps.PushNotifier = n
	return deps
}
