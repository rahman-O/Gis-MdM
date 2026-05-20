package groups

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	grouphttp "github.com/gis-mdm/server-backend-go/internal/modules/groups/adapter/http"
	grouppostgres "github.com/gis-mdm/server-backend-go/internal/modules/groups/adapter/persistence/postgres"
	groupapp "github.com/gis-mdm/server-backend-go/internal/modules/groups/application"
)

// Module registers device group routes.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "groups" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if deps.DB == nil {
		return fmt.Errorf("groups module requires DATABASE_URL")
	}
	repo := grouppostgres.NewGroupRepository(deps.DB)
	svc := groupapp.NewService(repo)
	grouphttp.NewHandler(svc).Register(groups.Private.Group("/groups"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
