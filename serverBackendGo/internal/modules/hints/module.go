package hints

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	hinthttp "github.com/gis-mdm/server-backend-go/internal/modules/hints/adapter/http"
	hintpostgres "github.com/gis-mdm/server-backend-go/internal/modules/hints/adapter/persistence/postgres"
	hintapp "github.com/gis-mdm/server-backend-go/internal/modules/hints/application"
)

// Module registers hints routes.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "hints" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if deps.DB == nil {
		return fmt.Errorf("hints module requires DATABASE_URL")
	}
	repo := hintpostgres.NewHintRepository(deps.DB)
	svc := hintapp.NewService(repo)
	hinthttp.NewHandler(svc).Register(groups.Private.Group("/hints"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
