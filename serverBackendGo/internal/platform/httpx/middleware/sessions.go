package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// SetupSessions installs cookie session store on the engine.
func SetupSessions(engine *gin.Engine, secret string) {
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: 2, // Lax
	})
	engine.Use(sessions.Sessions("hmdm_session", store))
}
