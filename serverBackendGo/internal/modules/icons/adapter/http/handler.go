package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	iconapp "github.com/gis-mdm/server-backend-go/internal/modules/icons/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/icons/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

type Handler struct {
	svc *iconapp.Service
}

func NewHandler(svc *iconapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/search", h.Search)
	g.GET("/search/:value", h.SearchByValue)
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

// Search godoc
// @Summary List icons
// @Tags Icons
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/icons/search [get]
func (h *Handler) Search(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	data, err := h.svc.Search(c.Request.Context(), p, "")
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, data)
}

// SearchByValue godoc
// @Summary Search icons by name
// @Tags Icons
// @Produce json
// @Security BearerAuth
// @Param value path string true "filter"
// @Success 200 {object} response.Envelope
// @Router /private/icons/search/{value} [get]
func (h *Handler) SearchByValue(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	data, err := h.svc.Search(c.Request.Context(), p, c.Param("value"))
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, data)
}

// Save godoc
// @Summary Create or update icon
// @Tags Icons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/icons [put]
func (h *Handler) Save(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var body domain.Icon
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	out, err := h.svc.Save(c.Request.Context(), p, body)
	if err != nil {
		if errors.Is(err, iconapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, out)
}

// Delete godoc
// @Summary Delete icon
// @Tags Icons
// @Produce json
// @Security BearerAuth
// @Param id path int true "icon id"
// @Success 200 {object} response.Envelope
// @Router /private/icons/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.Delete(c.Request.Context(), p, id); err != nil {
		if errors.Is(err, iconapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}
