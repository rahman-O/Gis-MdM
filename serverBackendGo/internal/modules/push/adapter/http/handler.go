package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	pushapp "github.com/gis-mdm/server-backend-go/internal/modules/push/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/push/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

type Handler struct {
	svc *pushapp.Service
}

func NewHandler(svc *pushapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.POST("", h.Send)
}

func (h *Handler) Send(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok || p == nil {
		response.PermissionDenied(c)
		return
	}
	var req domain.PushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "error.params.missing")
		return
	}
	if err := h.svc.Send(c.Request.Context(), p, req); err != nil {
		if errors.Is(err, pushapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}
