package roles

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	authpostgres "github.com/gis-mdm/server-backend-go/internal/modules/auth/adapter/persistence/postgres"
	roleshttp "github.com/gis-mdm/server-backend-go/internal/modules/roles/adapter/http"
	rolepostgres "github.com/gis-mdm/server-backend-go/internal/modules/roles/adapter/persistence/postgres"
	roleapp "github.com/gis-mdm/server-backend-go/internal/modules/roles/application"
)

// Module registers roles admin routes.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "roles" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if deps.DB == nil {
		return fmt.Errorf("roles module requires DATABASE_URL")
	}
	authRepo := authpostgres.NewUserRepository(deps.DB)
	repo := rolepostgres.NewRepository(deps.DB)
	svc := roleapp.NewService(repo, authRepo.IsSingleCustomer)
	h := roleshttp.NewHandler(svc)
	h.Register(groups.Private.Group("/roles"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
