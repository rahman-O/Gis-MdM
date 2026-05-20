package configurations

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	cfghttp "github.com/gis-mdm/server-backend-go/internal/modules/configurations/adapter/http"
	cfgpostgres "github.com/gis-mdm/server-backend-go/internal/modules/configurations/adapter/persistence/postgres"
)

// Module registers configuration list routes (Phase 4 subset).
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "configurations" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if deps.DB == nil {
		return fmt.Errorf("configurations module requires DATABASE_URL")
	}
	repo := cfgpostgres.NewConfigRepository(deps.DB)
	cfghttp.NewHandler(repo).Register(groups.Private.Group("/configurations"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
