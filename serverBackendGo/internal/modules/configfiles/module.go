package configfiles

import "github.com/gis-mdm/server-backend-go/internal/module"

// Module is a scaffold for gradual migration.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "configfiles" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	_ = groups.Private.Group("/config-files")
	deps.Log.Info("module scaffold registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
