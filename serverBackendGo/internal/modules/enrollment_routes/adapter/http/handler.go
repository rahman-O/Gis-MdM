package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	routeapp "github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

type Handler struct {
	svc *routeapp.Service
}

func NewHandler(svc *routeapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("", h.List)
	g.GET("/options/published-profile-versions", h.ListPublishedProfileVersions)
	g.POST("", h.Create)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", h.Update)
	g.GET("/:id/qr", h.QRMeta)
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
	case errors.Is(err, routeapp.ErrPermissionDenied):
		response.PermissionDenied(c)
	case errors.Is(err, routeapp.ErrRouteNotFound):
		response.ErrorEnvelope(c, "error.notfound.enrollment_route")
	case errors.Is(err, routeapp.ErrDuplicateRoute):
		response.DuplicateEntity(c, "error.duplicate.enrollment_route")
	case errors.Is(err, routeapp.ErrPublishedVersionRequired):
		response.ErrorEnvelope(c, "error.enrollment_route.published_version_required")
	case errors.Is(err, routeapp.ErrTreeNodeRequired):
		response.ErrorEnvelope(c, "error.enrollment_route.tree_node_required")
	case errors.Is(err, routeapp.ErrMainAppRequired):
		response.ErrorEnvelope(c, "error.enrollment_route.main_app_required")
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

func (h *Handler) ListPublishedProfileVersions(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	data, err := h.svc.ListPublishedProfileVersions(c.Request.Context(), p)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) GetByID(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	data, err := h.svc.GetByID(c.Request.Context(), p, id)
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
	var req domain.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "error.invalid.request")
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	var req domain.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	data, err := h.svc.Update(c.Request.Context(), p, id, req)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) QRMeta(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	data, err := h.svc.QRMeta(c.Request.Context(), p, id)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}
