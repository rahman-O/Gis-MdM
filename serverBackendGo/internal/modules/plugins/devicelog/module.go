package devicelog

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	dlhttp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/devicelog/adapter/http"
	dlpostgres "github.com/gis-mdm/server-backend-go/internal/modules/plugins/devicelog/adapter/persistence/postgres"
	dlapp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/devicelog/application"
)

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "plugins/devicelog" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModulePluginsEnabled || !deps.Config.ModulePluginsDevicelogEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("plugins/devicelog requires DATABASE_URL")
	}
	repo := dlpostgres.NewRepository(deps.DB)
	svc := dlapp.NewService(repo)
	dlhttp.NewHandler(svc).Register(groups.Plugins.Group("/devicelog"), groups.Engine)
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
