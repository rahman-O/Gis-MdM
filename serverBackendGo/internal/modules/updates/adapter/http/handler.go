package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	updapp "github.com/gis-mdm/server-backend-go/internal/modules/updates/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/updates/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

type Handler struct {
	svc *updapp.Service
}

func NewHandler(svc *updapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/check", h.Check)
	g.POST("", h.Apply)
}

func (h *Handler) Check(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	entries, err := h.svc.Check(c.Request.Context(), p)
	if err != nil {
		if errors.Is(err, updapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, entries)
}

func (h *Handler) Apply(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var req domain.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "error.params.missing")
		return
	}
	out, err := h.svc.Apply(c.Request.Context(), p, req)
	if err != nil {
		if errors.Is(err, updapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, out)
}
