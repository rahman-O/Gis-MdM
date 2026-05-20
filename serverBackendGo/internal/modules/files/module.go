package files

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	fileapp "github.com/gis-mdm/server-backend-go/internal/modules/files/application"
	filehttp "github.com/gis-mdm/server-backend-go/internal/modules/files/adapter/http"
	filepostgres "github.com/gis-mdm/server-backend-go/internal/modules/files/adapter/persistence/postgres"
	fileport "github.com/gis-mdm/server-backend-go/internal/modules/files/port"
	"github.com/gis-mdm/server-backend-go/internal/platform/storage"
)

// Module registers web UI file routes.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "files" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModuleFilesEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("files module requires DATABASE_URL")
	}
	store := storage.NewLocalStore(deps.Config.FilesDirectory)
	repo := filepostgres.NewFileRepository(deps.DB)
	customer := filepostgres.NewCustomerRepository(deps.DB)
	apps := filepostgres.NewAppLookup(deps.DB)
	svc := fileapp.NewService(repo, customer, apps, store, deps.Config.BaseURL, fileport.NoopPush())
	filehttp.NewHandler(svc).Register(groups.Private.Group("/web-ui-files"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
