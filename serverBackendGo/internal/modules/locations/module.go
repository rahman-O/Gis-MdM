package locations

import (
	"context"
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	lochttp "github.com/gis-mdm/server-backend-go/internal/modules/locations/adapter/http"
	locpostgres "github.com/gis-mdm/server-backend-go/internal/modules/locations/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/locations/application"
)

// Module registers location tracking routes and background workers.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "locations" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	cfg := deps.Config
	if !cfg.ModuleLocationsEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("locations module requires DATABASE_URL")
	}

	// Repository
	repo := locpostgres.NewRepository(deps.DB)

	// Application services
	dedup := application.NewDuplicateDetector(repo, deps.Log)
	limiter := application.NewLocationRateLimiter(
		cfg.LocationRateLimitMaxWrites,
		cfg.LocationRateLimitWindowSec,
		deps.Log,
	)
	broadcaster := lochttp.NewWebSocketServer(
		deps.Log,
		cfg.LocationWSHeartbeatSec,
		cfg.LocationWSDeviceTimeoutSec,
	)
	svc := application.NewLocationService(
		repo, dedup, limiter, broadcaster, deps.Log, cfg.LocationBatchMaxSize,
	)
	aggregator := application.NewAggregationWorker(
		repo, deps.Log,
		cfg.LocationRetentionDays,
		cfg.LocationAggregationIntervalH,
	)

	// REST handlers
	handler := lochttp.NewLocationHandler(svc, repo, deps.Log, cfg.LocationRetentionDays)
	sigMiddleware := lochttp.NewSignatureMiddleware(cfg, deps.Log)

	// Register REST routes
	api := groups.Engine.Group("/api/devices")
	api.POST("/:deviceId/locations/batch", sigMiddleware, handler.BatchUpload)
	api.GET("/:deviceId/locations", handler.GetHistory)
	api.GET("/:deviceId/locations/archive", handler.GetArchive)

	// Register WebSocket routes
	groups.Engine.GET("/ws/devices/:deviceId/location", broadcaster.HandleUpgrade)
	groups.Engine.GET("/ws/devices", broadcaster.HandleUpgrade)

	// Start background workers
	ctx := context.Background()
	go aggregator.Start(ctx)
	go limiter.StartCleanup(ctx)
	go broadcaster.StartHeartbeatMonitor(ctx)

	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
