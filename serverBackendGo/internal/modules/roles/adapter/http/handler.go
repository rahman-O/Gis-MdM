package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	roleapp "github.com/gis-mdm/server-backend-go/internal/modules/roles/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/roles/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves /private/roles/* endpoints.
type Handler struct {
	svc *roleapp.Service
}

func NewHandler(svc *roleapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/permissions", h.ListPermissions)
	g.GET("/all", h.ListAll)
	g.PUT("/", h.Upsert)
	g.DELETE("/:id", h.Delete)
}

func mapErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, roleapp.ErrPermissionDenied):
		response.PermissionDenied(c)
	case errors.Is(err, roleapp.ErrDuplicateRole):
		response.DuplicateEntity(c, "error.duplicate.role")
	default:
		response.ErrorEnvelope(c, "error.internal.server")
	}
}

// ListPermissions godoc
// @Summary List permissions
// @Tags Roles
// @Produce json
// @Success 200 {object} response.Envelope
// @Security BearerAuth
// @Router /private/roles/permissions [get]
func (h *Handler) ListPermissions(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	rows, err := h.svc.ListPermissions(c.Request.Context(), p)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, rows)
}

// ListAll godoc
// @Summary List all roles
// @Tags Roles
// @Produce json
// @Success 200 {object} response.Envelope
// @Security BearerAuth
// @Router /private/roles/all [get]
func (h *Handler) ListAll(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	rows, err := h.svc.ListRoles(c.Request.Context(), p)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, rows)
}

// Upsert godoc
// @Summary Create or update role
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/roles [put]
func (h *Handler) Upsert(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var body domain.RolePayload
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.UpsertRole(c.Request.Context(), p, body); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete godoc
// @Summary Delete role
// @Tags Roles
// @Param id path int true "Role ID"
// @Success 200 {object} response.Envelope
// @Security BearerAuth
// @Router /private/roles/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.DeleteRole(c.Request.Context(), p, id); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}
