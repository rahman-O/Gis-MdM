package audit

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	audithttp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/audit/adapter/http"
	auditpostgres "github.com/gis-mdm/server-backend-go/internal/modules/plugins/audit/adapter/persistence/postgres"
	auditapp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/audit/application"
)

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "plugins/audit" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if !deps.Config.ModulePluginsEnabled || !deps.Config.ModulePluginsAuditEnabled {
		deps.Log.Info("module disabled", "module", m.Name())
		return nil
	}
	if deps.DB == nil {
		return fmt.Errorf("plugins/audit requires DATABASE_URL")
	}
	repo := auditpostgres.NewAuditRepository(deps.DB)
	svc := auditapp.NewService(repo)
	audithttp.NewHandler(svc).Register(groups.Plugins.Group("/audit"))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
