package applications

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	appapp "github.com/gis-mdm/server-backend-go/internal/modules/applications/application"
	apphttp "github.com/gis-mdm/server-backend-go/internal/modules/applications/adapter/http"
	apppostgres "github.com/gis-mdm/server-backend-go/internal/modules/applications/adapter/persistence/postgres"
)

// Module registers application routes.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "applications" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if deps.DB == nil {
		return fmt.Errorf("applications module requires DATABASE_URL")
	}
	repo := apppostgres.NewApplicationRepository(deps.DB)
	svc := appapp.NewService(repo)
	apphttp.NewHandler(svc).Register(groups.Private.Group("/applications"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
