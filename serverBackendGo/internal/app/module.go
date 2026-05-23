package app

import (
	"github.com/gis-mdm/server-backend-go/internal/module"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx"
)

func toModuleGroups(g httpx.RouteGroups) module.RouteGroups {
	return module.RouteGroups{
		Engine:     g.Engine,
		Public:     g.Public,
		Private:    g.Private,
		Plugins:    g.Plugins,
		PluginMain:        g.PluginMain,
		PluginMainPrivate: g.PluginMainPrivate,
	}
}
