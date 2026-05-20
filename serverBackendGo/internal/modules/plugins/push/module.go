package push

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	pluginhttp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/adapter/http"
	pluginpostgres "github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/adapter/persistence/postgres"
	pluginapp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/application"
	notifpostgres "github.com/gis-mdm/server-backend-go/internal/modules/notifications/adapter/persistence/postgres"
)

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "plugins/push" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if deps.DB == nil {
		return fmt.Errorf("plugins/push module requires DATABASE_URL")
	}
	repo := pluginpostgres.NewMessageRepository(deps.DB)
	queue := notifpostgres.NewQueueRepository(deps.DB)
	svc := pluginapp.NewService(repo, queue)
	pluginhttp.NewHandler(svc).Register(groups.Plugins.Group("/push"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
