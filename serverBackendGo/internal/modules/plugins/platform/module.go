package platform

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	plathttp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/platform/adapter/http"
	platpostgres "github.com/gis-mdm/server-backend-go/internal/modules/plugins/platform/adapter/persistence/postgres"
	platapp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/platform/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/shared/status"
)

type Module struct {
	cache *status.Cache
}

func New() *Module {
	return &Module{cache: status.NewCache()}
}

func (m *Module) Name() string { return "plugins/platform" }

func (m *Module) Cache() *status.Cache { return m.cache }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModulePluginsEnabled || !deps.Config.ModulePluginsPlatformEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("plugins/platform requires DATABASE_URL")
	}
	repo := platpostgres.NewPluginRepository(deps.DB, deps.Config)
	svc := platapp.NewService(repo, m.cache)
	h := plathttp.NewHandler(svc)
	if groups.PluginMainPrivate != nil {
		h.RegisterPrivate(groups.PluginMainPrivate)
	}
	if groups.PluginMain != nil {
		h.RegisterPublic(groups.PluginMain.Group("/public"))
	}
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
