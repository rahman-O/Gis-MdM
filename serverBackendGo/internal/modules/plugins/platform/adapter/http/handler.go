package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	platapp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/platform/application"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

type Handler struct {
	svc *platapp.Service
}

func NewHandler(svc *platapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterPrivate(g *gin.RouterGroup) {
	g.GET("/available", h.Available)
	g.GET("/active", h.Active)
	g.POST("/disabled", h.Disabled)
}

func (h *Handler) RegisterPublic(g *gin.RouterGroup) {
	g.GET("/registered", h.Registered)
}

func (h *Handler) Available(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	out, err := h.svc.Available(c.Request.Context(), p)
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, out)
}

func (h *Handler) Active(c *gin.Context) {
	out, err := h.svc.Active(c.Request.Context())
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, out)
}

func (h *Handler) Registered(c *gin.Context) {
	out, err := h.svc.Registered(c.Request.Context())
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, out)
}

func (h *Handler) Disabled(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var ids []int64
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.ErrorEnvelope(c, "error.params.missing")
		return
	}
	if err := h.svc.SaveDisabled(c.Request.Context(), p, ids); err != nil {
		if errors.Is(err, platapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}
