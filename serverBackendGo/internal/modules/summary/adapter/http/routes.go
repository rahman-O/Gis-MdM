package http

import "github.com/gin-gonic/gin"

// Register mounts summary routes under /summary.
func Register(g *gin.RouterGroup, h *Handler) {
	g.GET("/devices", h.DeviceStats)
}
