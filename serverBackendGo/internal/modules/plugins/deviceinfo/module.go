package deviceinfo

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	dihttp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/deviceinfo/adapter/http"
	dipostgres "github.com/gis-mdm/server-backend-go/internal/modules/plugins/deviceinfo/adapter/persistence/postgres"
	diapp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/deviceinfo/application"
)

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "plugins/deviceinfo" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModulePluginsEnabled || !deps.Config.ModulePluginsDeviceinfoEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("plugins/deviceinfo requires DATABASE_URL")
	}
	repo := dipostgres.NewRepository(deps.DB)
	svc := diapp.NewService(repo)
	dihttp.NewHandler(svc).Register(groups.Plugins.Group("/deviceinfo"), groups.Engine)
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
