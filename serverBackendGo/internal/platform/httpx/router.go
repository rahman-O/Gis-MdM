package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/config"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	platformjwt "github.com/gis-mdm/server-backend-go/internal/platform/jwt"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/middleware"
)

// RouteGroups exposes legacy REST path prefixes for module registration.
type RouteGroups struct {
	Engine     *gin.Engine
	Public     *gin.RouterGroup
	Private    *gin.RouterGroup
	Plugins    *gin.RouterGroup
	PluginMain *gin.RouterGroup
}

// AuthWiring holds JWT and user lookup for protected routes.
type AuthWiring struct {
	JWT    *platformjwt.Provider
	Lookup platformauth.UserLookup
}

// BuildRouteGroups creates /rest/public, /rest/private, and plugin groups.
func BuildRouteGroups(engine *gin.Engine, cfg config.Config, wiring AuthWiring) RouteGroups {
	public := engine.Group("/rest/public")
	public.Use(middleware.IPFilter("PUBLIC_IP_ALLOWLIST"))

	jwtMW := middleware.JWTAuth(wiring.JWT, wiring.Lookup)

	private := engine.Group("/rest/private")
	private.Use(middleware.IPFilter("PRIVATE_IP_ALLOWLIST"))
	private.Use(jwtMW)
	private.Use(middleware.RequireAuth())

	plugins := engine.Group("/rest/plugins")
	plugins.Use(jwtMW)
	plugins.Use(middleware.RequireAuth())

	pluginMain := engine.Group("/rest/plugin/main")

	return RouteGroups{
		Engine:     engine,
		Public:     public,
		Private:    private,
		Plugins:    plugins,
		PluginMain: pluginMain,
	}
}
