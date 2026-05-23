package http

import (
	"github.com/gin-gonic/gin"
	statsapp "github.com/gis-mdm/server-backend-go/internal/modules/stats/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/stats/domain"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

type Handler struct {
	svc *statsapp.Service
}

func NewHandler(svc *statsapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.PUT("", h.PutStats)
}

// PutStats godoc
// @Summary Record server usage statistics
// @Tags Stats
// @Accept json
// @Produce json
// @Success 200 {object} response.Envelope
// @Router /public/stats [put]
func (h *Handler) PutStats(c *gin.Context) {
	var body domain.UsageStats
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.Save(c.Request.Context(), body); err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}
