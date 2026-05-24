package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	treeapp "github.com/gis-mdm/server-backend-go/internal/modules/device_tree/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/device_tree/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves /private/device-tree/* endpoints.
type Handler struct {
	svc *treeapp.Service
}

func NewHandler(svc *treeapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("", h.List)
	g.POST("/nodes", h.Create)
	g.PUT("/nodes/:id", h.Update)
	g.POST("/nodes/:id/delete", h.Delete)
}

func principal(c *gin.Context) (*platformauth.Principal, bool) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok || p == nil {
		c.Status(403)
		return nil, false
	}
	return p, true
}

func mapErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, treeapp.ErrPermissionDenied):
		response.PermissionDenied(c)
	case errors.Is(err, domain.ErrNotFound):
		response.ObjectNotFound(c)
	case errors.Is(err, domain.ErrDuplicateName):
		response.DuplicateEntity(c, "error.device_tree.duplicate_name")
	case errors.Is(err, domain.ErrInvalidParent):
		response.ErrorEnvelope(c, "error.device_tree.invalid_parent")
	case errors.Is(err, domain.ErrCycle):
		response.ErrorEnvelope(c, "error.device_tree.cycle")
	case errors.Is(err, domain.ErrTargetRequired):
		response.ErrorEnvelope(c, "error.device_tree.target_required")
	case errors.Is(err, domain.ErrCannotDeleteRoot):
		response.ErrorEnvelope(c, "error.device_tree.cannot_delete_root")
	default:
		response.ErrorEnvelope(c, "error.internal.server")
	}
}

func (h *Handler) List(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	data, err := h.svc.List(c.Request.Context(), p)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) Create(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var req domain.CreateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	data, err := h.svc.Create(c.Request.Context(), p, req)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) Update(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var req domain.UpdateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	data, err := h.svc.Update(c.Request.Context(), p, id, req)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) Delete(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var req domain.DeleteNodeRequest
	_ = c.ShouldBindJSON(&req)
	if err := h.svc.Delete(c.Request.Context(), p, id, req); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}
