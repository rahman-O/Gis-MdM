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
	svc := devapp.NewService(repo, port.NoopPush{})
	devhttp.NewHandler(svc).Register(groups.Private.Group("/devices"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
