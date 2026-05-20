package users

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	usershttp "github.com/gis-mdm/server-backend-go/internal/modules/users/adapter/http"
	userpostgres "github.com/gis-mdm/server-backend-go/internal/modules/users/adapter/persistence/postgres"
	userapp "github.com/gis-mdm/server-backend-go/internal/modules/users/application"
)

// Module registers user routes (profile + admin).
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "users" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if deps.DB == nil {
		return fmt.Errorf("users module requires DATABASE_URL")
	}
	repo := userpostgres.NewRepository(deps.DB)
	svc := userapp.NewService(repo)
	roles := userapp.NewRolesService(deps.DB)
	h := usershttp.NewHandler(svc, roles)
	h.Register(groups.Private.Group("/users"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
