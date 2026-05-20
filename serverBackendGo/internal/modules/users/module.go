package users

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/application"
	usershttp "github.com/gis-mdm/server-backend-go/internal/modules/users/adapter/http"
)

// Module registers user routes (current user for React shell).
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "users" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if deps.DB == nil {
		return fmt.Errorf("users module requires DATABASE_URL")
	}
	repo := postgres.NewUserRepository(deps.DB)
	svc := application.NewUsersService(repo)
	usershttp.Register(groups.Private.Group("/users"), usershttp.NewHandler(svc))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
