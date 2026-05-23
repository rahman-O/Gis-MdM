package notifications

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	notifhttp "github.com/gis-mdm/server-backend-go/internal/modules/notifications/adapter/http"
	notifpostgres "github.com/gis-mdm/server-backend-go/internal/modules/notifications/adapter/persistence/postgres"
	notifapp "github.com/gis-mdm/server-backend-go/internal/modules/notifications/application"
)

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "notifications" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModuleNotificationsEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("notifications module requires DATABASE_URL")
	}
	devRepo := notifpostgres.NewDeviceRepository(deps.DB)
	queue := notifpostgres.NewQueueRepository(deps.DB)
	svc := notifapp.NewService(devRepo, queue)
	h := notifhttp.NewHandler(svc, deps.Config.SecureEnrollment, deps.Config.HashSecret)
	h.Register(groups.Engine.Group("/rest/notifications"))
	h.RegisterPolling(groups.Engine)
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
