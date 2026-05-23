package passwordreset

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/application"
	"github.com/gis-mdm/server-backend-go/internal/platform/email"
	passwordresethttp "github.com/gis-mdm/server-backend-go/internal/modules/passwordreset/adapter/http"
)

// Module registers password reset routes.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "passwordreset" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModulePasswordResetEnabled {
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("passwordreset module requires DATABASE_URL")
	}
	repo := postgres.NewUserRepository(deps.DB)
	emailSvc := email.NewService(deps.Config.EmailConfigured, deps.Log)
	svc := application.NewPasswordResetService(repo, emailSvc)
	passwordresethttp.Register(groups.Public.Group("/passwordReset"), passwordresethttp.NewHandler(svc))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
