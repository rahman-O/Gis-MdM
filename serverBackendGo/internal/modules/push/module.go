package push

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	pushhttp "github.com/gis-mdm/server-backend-go/internal/modules/push/adapter/http"
	pushpostgres "github.com/gis-mdm/server-backend-go/internal/modules/push/adapter/persistence/postgres"
	pushapp "github.com/gis-mdm/server-backend-go/internal/modules/push/application"
	notifpostgres "github.com/gis-mdm/server-backend-go/internal/modules/notifications/adapter/persistence/postgres"
)

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "push" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModulePushEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("push module requires DATABASE_URL")
	}
	targets := pushpostgres.NewTargetRepository(deps.DB)
	queue := notifpostgres.NewQueueRepository(deps.DB)
	svc := pushapp.NewService(targets, queue)
	pushhttp.NewHandler(svc).Register(groups.Private.Group("/push"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
