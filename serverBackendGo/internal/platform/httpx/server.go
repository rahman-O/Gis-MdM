package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/config"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/middleware"
)

// NewEngine builds the Gin engine with global middleware.
func NewEngine(cfg config.Config) *gin.Engine {
	gin.SetMode(cfg.GinMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())
	return r
}
