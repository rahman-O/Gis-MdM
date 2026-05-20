package groups

import "github.com/gis-mdm/server-backend-go/internal/module"

// Module is a scaffold for gradual migration.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "groups" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	_ = groups.Private.Group("/groups")
	deps.Log.Info("module scaffold registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
