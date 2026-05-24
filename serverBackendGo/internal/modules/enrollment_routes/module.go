package enrollment_routes

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	routehttp "github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/adapter/http"
	routepostgres "github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/adapter/persistence/postgres"
	routeapp "github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/application"
)

// Module registers enrollment route routes (017 US4).
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "enrollment_routes" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModuleEnrollmentRoutesEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("enrollment_routes module requires DATABASE_URL")
	}
	repo := routepostgres.NewRouteRepository(deps.DB)
	svc := routeapp.NewService(repo)
	routehttp.NewHandler(svc).Register(groups.Private.Group("/enrollment-routes"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
