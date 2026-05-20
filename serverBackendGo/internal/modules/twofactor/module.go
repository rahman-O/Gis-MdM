package twofactor

import (
	"github.com/gis-mdm/server-backend-go/internal/module"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/application"
	twofactorhttp "github.com/gis-mdm/server-backend-go/internal/modules/twofactor/adapter/http"
)

// Module registers two-factor HTTP routes (requires auth session/JWT).
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "twofactor" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if deps.DB == nil {
		return nil
	}
	repo := postgres.NewUserRepository(deps.DB)
	svc := application.NewTwoFactorService(repo, "Headwind MDM")
	twofactorhttp.Register(groups.Private.Group("/twofactor"), twofactorhttp.NewHandler(svc))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
