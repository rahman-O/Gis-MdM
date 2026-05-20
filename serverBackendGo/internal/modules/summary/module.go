package summary

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	summaryhttp "github.com/gis-mdm/server-backend-go/internal/modules/summary/adapter/http"
	"github.com/gis-mdm/server-backend-go/internal/modules/summary/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/summary/application"
)

// Module registers dashboard summary routes.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "summary" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if deps.DB == nil {
		return fmt.Errorf("summary module requires DATABASE_URL")
	}
	repo := postgres.NewSummaryRepository(deps.DB)
	svc := application.NewService(repo)
	summaryhttp.Register(groups.Private.Group("/summary"), summaryhttp.NewHandler(svc))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
