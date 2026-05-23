package icons

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	iconapp "github.com/gis-mdm/server-backend-go/internal/modules/icons/application"
	iconhttp "github.com/gis-mdm/server-backend-go/internal/modules/icons/adapter/http"
	iconpostgres "github.com/gis-mdm/server-backend-go/internal/modules/icons/adapter/persistence/postgres"
)

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "icons" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModuleIconsEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("icons module requires DATABASE_URL")
	}
	repo := iconpostgres.NewIconRepository(deps.DB)
	svc := iconapp.NewService(repo)
	iconhttp.NewHandler(svc).Register(groups.Private.Group("/icons"))
	uploadRepo := iconpostgres.NewUploadedFileRepository(deps.DB)
	store := iconapp.IconFileStore{FilesDir: deps.Config.FilesDirectory}
	iconhttp.NewIconFileHandler(store, uploadRepo, uploadRepo).Register(groups.Private.Group("/icon-files"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
