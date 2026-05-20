package signup

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/application"
	"github.com/gis-mdm/server-backend-go/internal/platform/email"
	signuphttp "github.com/gis-mdm/server-backend-go/internal/modules/signup/adapter/http"
)

// Module registers customer signup routes.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "signup" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModuleSignupEnabled {
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("signup module requires DATABASE_URL")
	}
	repo := postgres.NewUserRepository(deps.DB)
	emailSvc := email.NewService(deps.Config.EmailConfigured, deps.Log)
	svc := application.NewSignupService(repo, emailSvc, deps.Config.CustomerSignup)
	signuphttp.Register(groups.Public.Group("/signup"), signuphttp.NewHandler(svc))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
