package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	authdomain "github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	userapp "github.com/gis-mdm/server-backend-go/internal/modules/users/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/users/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves /private/users/* endpoints.
type Handler struct {
	svc   *userapp.Service
	roles *userapp.RolesService
}

func NewHandler(svc *userapp.Service, roles *userapp.RolesService) *Handler {
	return &Handler{svc: svc, roles: roles}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/current", h.Current)
	g.PUT("/details", h.UpdateDetails)
	g.PUT("/current", h.ChangePassword)
	g.GET("/all", h.ListAll)
	g.PUT("/", h.Upsert)
	g.DELETE("/other/:id", h.DeleteOther)
	g.GET("/roles", h.ListRoles)
}

func principal(c *gin.Context) (*platformauth.Principal, bool) {
	return platformauth.PrincipalFromContext(c)
}

func mapSvcErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, userapp.ErrPermissionDenied):
		response.PermissionDenied(c)
	case errors.Is(err, userapp.ErrWrongPassword):
		response.ErrorEnvelope(c, "error.password.wrong")
	case errors.Is(err, userapp.ErrEmptyPassword):
		response.ErrorEnvelope(c, "error.password.empty")
	case errors.Is(err, userapp.ErrDuplicateEmail):
		response.ErrorEnvelope(c, "error.duplicate.email")
	case errors.Is(err, userapp.ErrDuplicateLogin):
		response.ErrorEnvelope(c, "error.duplicate.login")
	case errors.Is(err, userapp.ErrUserNotFound):
		response.ErrorEnvelope(c, "error.user.not.found")
	default:
		response.ErrorEnvelope(c, "error.internal.server")
	}
}

// Current godoc
// @Summary Get current user
// @Tags Users
// @Produce json
// @Success 200 {object} response.Envelope
// @Security BearerAuth
// @Router /private/users/current [get]
func (h *Handler) Current(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		response.OK(c, nil)
		return
	}
	view, err := h.svc.GetCurrentUser(c.Request.Context(), p.ID)
	if err != nil {
		mapSvcErr(c, err)
		return
	}
	response.OK(c, view)
}

// UpdateDetails godoc
// @Summary Update profile details
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} response.Envelope
// @Security BearerAuth
// @Router /private/users/details [put]
func (h *Handler) UpdateDetails(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var body domain.ProfilePayload
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	view, err := h.svc.UpdateProfile(c.Request.Context(), p, body)
	if err != nil {
		mapSvcErr(c, err)
		return
	}
	response.OKMessage(c, "success.operation.completed", view)
}

// ChangePassword godoc
// @Summary Change current user password
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} response.Envelope
// @Security BearerAuth
// @Router /private/users/current [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var body domain.UserPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if body.ID == nil {
		id := p.ID
		body.ID = &id
	}
	if err := h.svc.ChangePassword(c.Request.Context(), p, body); err != nil {
		mapSvcErr(c, err)
		return
	}
	response.OKMessage(c, "success.operation.completed", nil)
}

// ListAll godoc
// @Summary List tenant users
// @Tags Users
// @Produce json
// @Success 200 {object} response.Envelope
// @Security BearerAuth
// @Router /private/users/all [get]
func (h *Handler) ListAll(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	rows, err := h.svc.ListUsers(c.Request.Context(), p, c.Query("filter"))
	if err != nil {
		mapSvcErr(c, err)
		return
	}
	if rows == nil {
		rows = []*authdomain.UserView{}
	}
	response.OK(c, rows)
}

// Upsert godoc
// @Summary Create or update user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/users [put]
func (h *Handler) Upsert(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var body domain.UserPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.UpsertUser(c.Request.Context(), p, body); err != nil {
		mapSvcErr(c, err)
		return
	}
	response.OK(c, nil)
}

// DeleteOther godoc
// @Summary Delete user by id
// @Tags Users
// @Param id path int true "User ID"
// @Success 200 {object} response.Envelope
// @Security BearerAuth
// @Router /private/users/other/{id} [delete]
func (h *Handler) DeleteOther(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.DeleteUser(c.Request.Context(), p, id); err != nil {
		mapSvcErr(c, err)
		return
	}
	response.OK(c, nil)
}

// ListRoles godoc
// @Summary List assignable roles
// @Tags Users
// @Produce json
// @Success 200 {object} response.Envelope
// @Security BearerAuth
// @Router /private/users/roles [get]
func (h *Handler) ListRoles(c *gin.Context) {
	if h.roles == nil {
		response.OK(c, []userapp.RoleRow{})
		return
	}
	rows, err := h.roles.ListRoles(c.Request.Context())
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	if rows == nil {
		rows = []userapp.RoleRow{}
	}
	response.OK(c, rows)
}
