package device_tree

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	treehttp "github.com/gis-mdm/server-backend-go/internal/modules/device_tree/adapter/http"
	treepostgres "github.com/gis-mdm/server-backend-go/internal/modules/device_tree/adapter/persistence/postgres"
	treeapp "github.com/gis-mdm/server-backend-go/internal/modules/device_tree/application"
)

// Module registers device tree routes.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "device_tree" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModuleDeviceTreeEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("device_tree module requires DATABASE_URL")
	}
	repo := treepostgres.NewTreeRepository(deps.DB)
	svc := treeapp.NewService(repo)
	treehttp.NewHandler(svc).Register(groups.Private.Group("/device-tree"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
