package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	groupapp "github.com/gis-mdm/server-backend-go/internal/modules/groups/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/groups/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves /private/groups/* endpoints.
type Handler struct {
	svc *groupapp.Service
}

func NewHandler(svc *groupapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/search", h.Search)
	g.GET("/search/:value", h.SearchByValue)
	g.POST("/autocomplete", h.Autocomplete)
	g.PUT("", h.Save)
	g.DELETE("/:id", h.Delete)
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
	case errors.Is(err, groupapp.ErrPermissionDenied):
		response.PermissionDenied(c)
	case errors.Is(err, groupapp.ErrDuplicateGroup):
		response.DuplicateEntity(c, "error.duplicate.group")
	case errors.Is(err, groupapp.ErrNotEmptyGroup):
		response.ErrorEnvelope(c, "error.notempty.group")
	default:
		response.ErrorEnvelope(c, "error.internal.server")
	}
}

// Search godoc
// @Summary List device groups
// @Tags Groups
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/groups/search [get]
func (h *Handler) Search(c *gin.Context) {
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

// SearchByValue godoc
// @Summary Search device groups by name
// @Tags Groups
// @Produce json
// @Security BearerAuth
// @Param value path string true "Filter"
// @Success 200 {object} response.Envelope
// @Router /private/groups/search/{value} [get]
func (h *Handler) SearchByValue(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	data, err := h.svc.SearchByValue(c.Request.Context(), p, c.Param("value"))
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// Autocomplete godoc
// @Summary Group autocomplete
// @Tags Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/groups/autocomplete [post]
func (h *Handler) Autocomplete(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var filter string
	_ = c.ShouldBindJSON(&filter)
	data, err := h.svc.Autocomplete(c.Request.Context(), p, filter)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// Save godoc
// @Summary Create or update device group
// @Tags Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/groups [put]
func (h *Handler) Save(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var g domain.Group
	if err := c.ShouldBindJSON(&g); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.Save(c.Request.Context(), p, g); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete godoc
// @Summary Delete device group
// @Tags Groups
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Success 200 {object} response.Envelope
// @Router /private/groups/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.Delete(c.Request.Context(), p, id); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}
