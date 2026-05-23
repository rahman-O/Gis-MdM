package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	msgapp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/messaging/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/messaging/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

type Handler struct {
	svc *msgapp.Service
}

func NewHandler(svc *msgapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.POST("/private/search", h.Search)
	g.POST("/private/send", h.Send)
	g.DELETE("/:id", h.Delete)
	g.GET("/private/purge/:days", h.Purge)
}

func (h *Handler) RegisterPublic(engine *gin.Engine) {
	engine.GET("/rest/plugins/messaging/public/status/:id/:status", h.Status)
}

func (h *Handler) Search(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var f domain.MessageFilter
	_ = c.ShouldBindJSON(&f)
	out, err := h.svc.Search(c.Request.Context(), p, f)
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, out)
}

func (h *Handler) Send(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var req domain.SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "error.params.missing")
		return
	}
	if err := h.svc.Send(c.Request.Context(), p, req); err != nil {
		if errors.Is(err, msgapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) Delete(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), p, id); err != nil {
		if errors.Is(err, msgapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

func (h *Handler) Purge(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	days, _ := strconv.Atoi(c.Param("days"))
	if err := h.svc.Purge(c.Request.Context(), p, days); err != nil {
		if errors.Is(err, msgapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

func (h *Handler) Status(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	status, _ := strconv.Atoi(c.Param("status"))
	if err := h.svc.SetStatus(c.Request.Context(), id, status); err != nil {
		response.ErrorEnvelope(c, err.Error())
		return
	}
	response.OK(c, nil)
}
