package http

import "github.com/gin-gonic/gin"

// Register mounts auth routes on the public REST group.
func Register(public *gin.RouterGroup, h *Handler) {
	auth := public.Group("/auth")
	auth.GET("/options", h.Options)
	auth.POST("/login", h.Login)
	auth.POST("/logout", h.Logout)
	public.POST("/jwt/login", h.JWTLogin)
}
