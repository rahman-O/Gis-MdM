package configfiles

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	cfgfilehttp "github.com/gis-mdm/server-backend-go/internal/modules/configfiles/adapter/http"
	cfgfilepostgres "github.com/gis-mdm/server-backend-go/internal/modules/configfiles/adapter/persistence/postgres"
)

// Module registers config file upload routes.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "configfiles" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if deps.DB == nil {
		return fmt.Errorf("configfiles module requires DATABASE_URL")
	}
	h := cfgfilehttp.NewHandler(
		deps.Config.FilesDirectory,
		deps.Config.BaseURL,
		cfgfilepostgres.NewCustomerFilesRepo(deps.DB),
	)
	h.Register(groups.Private.Group("/config-files"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
