package deviceinfo

import "github.com/gis-mdm/server-backend-go/internal/module"

type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "plugins/deviceinfo" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
	_ = groups.Plugins
	deps.Log.Info("module scaffold registered", "module", m.Name())
	return nil
}

var _ module.Module = (*Module)(nil)
