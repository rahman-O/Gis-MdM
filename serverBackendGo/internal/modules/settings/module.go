package settings

import (
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/module"
	settingshttp "github.com/gis-mdm/server-backend-go/internal/modules/settings/adapter/http"
	"github.com/gis-mdm/server-backend-go/internal/modules/settings/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/settings/application"
)

// Module registers tenant settings routes.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "settings" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	if deps.DB == nil {
		return fmt.Errorf("settings module requires DATABASE_URL")
	}
	repo := postgres.NewSettingsRepository(deps.DB)
	svc := application.NewService(repo)
	settingshttp.Register(groups.Private.Group("/settings"), settingshttp.NewHandler(svc))
	deps.Log.Info("module registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
