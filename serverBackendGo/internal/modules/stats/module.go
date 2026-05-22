package stats

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	statshttp "github.com/gis-mdm/server-backend-go/internal/modules/stats/adapter/http"
	statspostgres "github.com/gis-mdm/server-backend-go/internal/modules/stats/adapter/persistence/postgres"
	statsapp "github.com/gis-mdm/server-backend-go/internal/modules/stats/application"
)

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "stats" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModuleStatsEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("stats module requires DATABASE_URL")
	}
	repo := statspostgres.NewStatsRepository(deps.DB)
	svc := statsapp.NewService(repo)
	statshttp.NewHandler(svc).Register(groups.Public.Group("/stats"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
