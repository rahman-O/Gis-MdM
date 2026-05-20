package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	custapp "github.com/gis-mdm/server-backend-go/internal/modules/customers/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/customers/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves /private/customers/* endpoints.
type Handler struct {
	svc *custapp.Service
}

func NewHandler(svc *custapp.Service) *Handler {
	return &Handler{svc: svc}
}

// Register mounts routes on /customers.
func (h *Handler) Register(g *gin.RouterGroup) {
	g.POST("/search", h.Search)
	g.GET("/impersonate/:id", h.Impersonate)
	g.GET("/prefix/:prefix/used", h.PrefixUsed)
	g.GET("/:id/edit", h.GetForEdit)
	g.DELETE("/:id", h.Delete)
	g.PUT("", h.Save)
}

func mapErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, custapp.ErrPermissionDenied):
		response.PermissionDenied(c)
	case errors.Is(err, custapp.ErrDuplicateCustomerName):
		response.DuplicateEntity(c, "error.duplicate.customer.name")
	case errors.Is(err, custapp.ErrDuplicateEmail):
		response.DuplicateEntity(c, "error.duplicate.email")
	case errors.Is(err, custapp.ErrOrgAdminNotFound):
		response.ErrorEnvelope(c, "error.notfound.customer.admin")
	case errors.Is(err, custapp.ErrImpersonationBlocked):
		response.ErrorEnvelope(c, "error.permission.denied")
	default:
		response.ErrorEnvelope(c, "error.internal.server")
	}
}

// Search godoc
// @Summary Search customers
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/customers/search [post]
func (h *Handler) Search(c *gin.Context) {
	p, ok := platformauth.RequireSuperAdmin(c)
	if !ok {
		return
	}
	var req domain.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	data, err := h.svc.Search(c.Request.Context(), p, req)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// Impersonate godoc
// @Summary Impersonate customer org admin
// @Tags Customers
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/customers/impersonate/{id} [get]
func (h *Handler) Impersonate(c *gin.Context) {
	p, ok := platformauth.RequireSuperAdmin(c)
	if !ok {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	view, err := h.svc.Impersonate(c.Request.Context(), p, id)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, view)
}

// Save godoc
// @Summary Create or update customer
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/customers [put]
func (h *Handler) Save(c *gin.Context) {
	p, ok := platformauth.RequireSuperAdmin(c)
	if !ok {
		return
	}
	var body domain.Customer
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	data, err := h.svc.Save(c.Request.Context(), p, body)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// Delete godoc
// @Summary Delete customer
// @Tags Customers
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/customers/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	p, ok := platformauth.RequireSuperAdmin(c)
	if !ok {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), p, id); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}

// GetForEdit godoc
// @Summary Get customer for edit
// @Tags Customers
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/customers/{id}/edit [get]
func (h *Handler) GetForEdit(c *gin.Context) {
	p, ok := platformauth.RequireSuperAdmin(c)
	if !ok {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	cust, err := h.svc.GetForEdit(c.Request.Context(), p, id)
	if err != nil {
		mapErr(c, err)
		return
	}
	if cust == nil {
		response.ObjectNotFound(c)
		return
	}
	response.OK(c, cust)
}

// PrefixUsed godoc
// @Summary Check if device prefix is used
// @Tags Customers
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/customers/prefix/{prefix}/used [get]
func (h *Handler) PrefixUsed(c *gin.Context) {
	p, ok := platformauth.RequireSuperAdmin(c)
	if !ok {
		return
	}
	used, err := h.svc.PrefixUsed(c.Request.Context(), p, c.Param("prefix"))
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, used)
}
