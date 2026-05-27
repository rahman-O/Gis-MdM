package devices

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	devhttp "github.com/gis-mdm/server-backend-go/internal/modules/devices/adapter/http"
	devpostgres "github.com/gis-mdm/server-backend-go/internal/modules/devices/adapter/persistence/postgres"
	devapp "github.com/gis-mdm/server-backend-go/internal/modules/devices/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/devices/port"
)

// Module registers device routes.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "devices" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if deps.DB == nil {
		return fmt.Errorf("devices module requires DATABASE_URL")
	}
	repo := devpostgres.NewDeviceRepository(deps.DB)
	var push port.PushNotifier = port.NoopPush{}
	if deps.PushNotifier != nil {
		if n, ok := deps.PushNotifier.(port.PushNotifier); ok {
			push = n
		}
	}
	svc := devapp.NewService(repo, push)
	devhttp.NewHandler(svc).Register(groups.Private.Group("/devices"))

	// Location tracking endpoints
	locHandler := devhttp.NewLocationHandler(deps.DB)
	locHandler.RegisterPublic(groups.Public)
	locHandler.RegisterPrivate(groups.Private)

	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
