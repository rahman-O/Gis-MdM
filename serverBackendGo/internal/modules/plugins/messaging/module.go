package messaging

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	msghttp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/messaging/adapter/http"
	msgpostgres "github.com/gis-mdm/server-backend-go/internal/modules/plugins/messaging/adapter/persistence/postgres"
	msgapp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/messaging/application"
	notifpostgres "github.com/gis-mdm/server-backend-go/internal/modules/notifications/adapter/persistence/postgres"
	targetpostgres "github.com/gis-mdm/server-backend-go/internal/modules/plugins/shared/targets/postgres"
)

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "plugins/messaging" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModulePluginsEnabled || !deps.Config.ModulePluginsMessagingEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("plugins/messaging requires DATABASE_URL")
	}
	repo := msgpostgres.NewMessageRepository(deps.DB)
	queue := notifpostgres.NewQueueRepository(deps.DB)
	targets := targetpostgres.NewResolver(deps.DB)
	svc := msgapp.NewService(repo, queue, targets)
	h := msghttp.NewHandler(svc)
	h.Register(groups.Plugins.Group("/messaging"))
	h.RegisterPublic(groups.Engine)
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
