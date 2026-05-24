package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	profileapp "github.com/gis-mdm/server-backend-go/internal/modules/profiles/application"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// RegisterHubRoutes adds 019 summary/activity routes (register before /:id catch-alls if needed).
func RegisterHubRoutes(g *gin.RouterGroup, hub *profileapp.HubService) {
	if hub == nil {
		return
	}
	g.GET("/:id/summary", func(c *gin.Context) { hubSummary(c, hub) })
	g.GET("/:id/activity", func(c *gin.Context) { hubActivity(c, hub) })
}

func hubSummary(c *gin.Context, hub *profileapp.HubService) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, err := strconv.Atoi(c.Param("id"))
	if err != nil || profileID <= 0 {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	data, err := hub.Summary(c.Request.Context(), p, profileID)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

func hubActivity(c *gin.Context, hub *profileapp.HubService) {
	p, ok := principal(c)
	if !ok {
		return
	}
	profileID, err := strconv.Atoi(c.Param("id"))
	if err != nil || profileID <= 0 {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	limit := 50
	if q := c.Query("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 {
			limit = n
		}
	}
	data, err := hub.Activity(c.Request.Context(), p, profileID, limit)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}
