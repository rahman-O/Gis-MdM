package publicapi

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	pubapp "github.com/gis-mdm/server-backend-go/internal/modules/publicapi/application"
	pubhttp "github.com/gis-mdm/server-backend-go/internal/modules/publicapi/adapter/http"
	pubpostgres "github.com/gis-mdm/server-backend-go/internal/modules/publicapi/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/platform/storage"
)

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "publicapi" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModulePublicAPIEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("publicapi module requires DATABASE_URL")
	}
	repo := pubpostgres.NewDeviceRepository(deps.DB)
	store := storage.NewLocalStore(deps.Config.FilesDirectory)
	cfg := pubapp.RebrandingConfig{
		AppName:     deps.Config.RebrandingName,
		VendorName:  deps.Config.RebrandingVendor,
		VendorLink:  deps.Config.RebrandingVendorURL,
		SignupLink:  deps.Config.RebrandingSignupURL,
		TermsLink:   deps.Config.RebrandingTermsURL,
		LogoPath:    deps.Config.RebrandingLogo,
		HashSecret:  deps.Config.HashSecret,
	}
	svc := pubapp.NewService(repo, store, deps.Config.BaseURL, cfg)
	pubhttp.NewHandler(svc).Register(groups.Public)
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
