package sync

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	synchttp "github.com/gis-mdm/server-backend-go/internal/modules/sync/adapter/http"
	syncpostgres "github.com/gis-mdm/server-backend-go/internal/modules/sync/adapter/persistence/postgres"
	syncapp "github.com/gis-mdm/server-backend-go/internal/modules/sync/application"
)

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "sync" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModuleSyncEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("sync module requires DATABASE_URL")
	}
	repo := syncpostgres.NewDeviceSyncRepository(deps.DB)
	svc := syncapp.NewService(repo, syncapp.Config{
		BaseURL:           deps.Config.BaseURL,
		FilesDirectory:    deps.Config.FilesDirectory,
		HashSecret:        deps.Config.HashSecret,
		SecureEnrollment:  deps.Config.SecureEnrollment,
		PreventDuplicate:  deps.Config.PreventDuplicateEnrollment,
		MobileAppName:     deps.Config.RebrandingMobileName,
		VendorName:        deps.Config.RebrandingVendor,
		DefaultCustomerID: 1,
	})
	synchttp.NewHandler(svc, deps.Config.HashSecret).Register(groups.Public.Group("/sync"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
