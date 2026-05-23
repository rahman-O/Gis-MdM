package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/modules/summary/application"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves summary endpoints.
type Handler struct {
	svc *application.Service
}

func NewHandler(svc *application.Service) *Handler {
	return &Handler{svc: svc}
}

// DeviceStats godoc
// @Summary Device dashboard statistics
// @Tags Summary
// @Produce json
// @Success 200 {object} response.Envelope
// @Security BearerAuth
// @Router /private/summary/devices [get]
func (h *Handler) DeviceStats(c *gin.Context) {
	principal, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.ErrorEnvelope(c, "error.permission.denied")
		return
	}
	stats, err := h.svc.GetDeviceStats(c.Request.Context(), principal.CustomerID, principal.ID)
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, stats)
}
