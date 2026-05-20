package customers

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	custhttp "github.com/gis-mdm/server-backend-go/internal/modules/customers/adapter/http"
	custpostgres "github.com/gis-mdm/server-backend-go/internal/modules/customers/adapter/persistence/postgres"
	custapp "github.com/gis-mdm/server-backend-go/internal/modules/customers/application"
)

// Module registers customers routes.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "customers" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if deps.DB == nil {
		return fmt.Errorf("customers module requires DATABASE_URL")
	}
	repo := custpostgres.NewCustomerRepository(deps.DB)
	users := custpostgres.NewUserLookup(deps.DB)
	svc := custapp.NewService(repo, users)
	custhttp.NewHandler(svc).Register(groups.Private.Group("/customers"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
