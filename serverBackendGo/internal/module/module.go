package module

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/config"
)

// Module is implemented by each domain package under internal/modules.
type Module interface {
	Name() string
	Register(groups RouteGroups, deps Dependencies) error
}

// RouteGroups mirrors legacy REST prefixes.
type RouteGroups struct {
	Engine     *gin.Engine
	Public     *gin.RouterGroup
	Private    *gin.RouterGroup
	Plugins           *gin.RouterGroup
	PluginMain        *gin.RouterGroup
	PluginMainPrivate *gin.RouterGroup
}

// Dependencies shared across modules during registration.
type Dependencies struct {
	Config config.Config
	DB     *sql.DB
	Log    *slog.Logger
	// PushNotifier enqueues configUpdated/appConfigUpdated (nil → modules use noop).
	PushNotifier interface {
		NotifyConfigurationChanged(int) error
		NotifyAppSettings(context.Context, int) error
	}
}
