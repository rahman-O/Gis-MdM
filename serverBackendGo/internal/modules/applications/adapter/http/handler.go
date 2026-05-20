package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	appapp "github.com/gis-mdm/server-backend-go/internal/modules/applications/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/applications/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves /private/applications/* endpoints.
type Handler struct {
	svc *appapp.Service
}

func NewHandler(svc *appapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/search", h.Search)
	g.GET("/search/:value", h.SearchByValue)
	g.POST("/autocomplete", h.Autocomplete)
	g.GET("/admin/search", h.AdminSearch)
	g.GET("/admin/search/:value", h.AdminSearchByValue)
	g.GET("/admin/common/:id", h.TurnIntoCommon)
	g.PUT("/android", h.SaveAndroid)
	g.PUT("/web", h.SaveWeb)
	g.PUT("/versions", h.SaveVersion)
	g.PUT("/validatePkg", h.ValidatePkg)
	g.GET("/configurations/:id", h.GetAppConfigurations)
	g.POST("/configurations", h.UpdateAppConfigurations)
	g.GET("/version/:versionId/configurations", h.GetVersionConfigurations)
	g.POST("/version/configurations", h.UpdateVersionConfigurations)
	g.GET("/:id/versions", h.ListVersions)
	g.GET("/:id", h.GetByID)
	g.DELETE("/versions/:id", h.DeleteVersion)
	g.DELETE("/:id", h.DeleteApp)
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
	case errors.Is(err, appapp.ErrPermissionDenied):
		response.PermissionDenied(c)
	case errors.Is(err, appapp.ErrAppNotFound):
		response.ErrorEnvelope(c, "error.notfound.application")
	default:
		response.ErrorEnvelope(c, "error.internal.server")
	}
}

// Search godoc
// @Summary List applications
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/applications/search [get]
func (h *Handler) Search(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	data, err := h.svc.Search(c.Request.Context(), p)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// SearchByValue godoc
// @Summary Search applications by name or package
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Param value path string true "Filter"
// @Success 200 {object} response.Envelope
// @Router /private/applications/search/{value} [get]
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
// @Summary Application autocomplete
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/applications/autocomplete [post]
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

// GetByID godoc
// @Summary Get application
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Success 200 {object} response.Envelope
// @Router /private/applications/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.svc.GetByID(c.Request.Context(), p, id)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// ListVersions godoc
// @Summary List application versions
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Success 200 {object} response.Envelope
// @Router /private/applications/{id}/versions [get]
func (h *Handler) ListVersions(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.svc.ListVersions(c.Request.Context(), p, id)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// SaveAndroid godoc
// @Summary Create or update Android application
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/applications/android [put]
func (h *Handler) SaveAndroid(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var app domain.Application
	if err := c.ShouldBindJSON(&app); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	data, err := h.svc.SaveAndroid(c.Request.Context(), p, app)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// SaveWeb godoc
// @Summary Create or update web application
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/applications/web [put]
func (h *Handler) SaveWeb(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var app domain.Application
	if err := c.ShouldBindJSON(&app); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	data, err := h.svc.SaveWeb(c.Request.Context(), p, app)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// SaveVersion godoc
// @Summary Create or update application version
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/applications/versions [put]
func (h *Handler) SaveVersion(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var ver domain.ApplicationVersion
	if err := c.ShouldBindJSON(&ver); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	data, err := h.svc.SaveVersion(c.Request.Context(), p, ver)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// ValidatePkg godoc
// @Summary Validate package id uniqueness
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/applications/validatePkg [put]
func (h *Handler) ValidatePkg(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var req domain.ValidatePkgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	data, err := h.svc.ValidatePkg(c.Request.Context(), p, req)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// DeleteApp godoc
// @Summary Delete application
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Success 200 {object} response.Envelope
// @Router /private/applications/{id} [delete]
func (h *Handler) DeleteApp(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.DeleteApp(c.Request.Context(), p, id); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}

// DeleteVersion godoc
// @Summary Delete application version
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Param id path int true "Version ID"
// @Success 200 {object} response.Envelope
// @Router /private/applications/versions/{id} [delete]
func (h *Handler) DeleteVersion(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.DeleteVersion(c.Request.Context(), p, id); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}

// GetAppConfigurations godoc
// @Summary Configurations linked to application
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Success 200 {object} response.Envelope
// @Router /private/applications/configurations/{id} [get]
func (h *Handler) GetAppConfigurations(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.svc.GetAppConfigurations(c.Request.Context(), p, id)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// UpdateAppConfigurations godoc
// @Summary Update configuration links for application
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/applications/configurations [post]
func (h *Handler) UpdateAppConfigurations(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var req domain.LinkConfigurationsToAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.UpdateAppConfigurations(c.Request.Context(), p, req); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}

// GetVersionConfigurations godoc
// @Summary Configurations linked to application version
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Param versionId path int true "Version ID"
// @Success 200 {object} response.Envelope
// @Router /private/applications/version/{versionId}/configurations [get]
func (h *Handler) GetVersionConfigurations(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("versionId"))
	data, err := h.svc.GetVersionConfigurations(c.Request.Context(), p, id)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// UpdateVersionConfigurations godoc
// @Summary Update configuration links for application version
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/applications/version/configurations [post]
func (h *Handler) UpdateVersionConfigurations(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var req domain.LinkConfigurationsToAppVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.UpdateVersionConfigurations(c.Request.Context(), p, req); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}

// AdminSearch godoc
// @Summary List shared applications (super-admin)
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/applications/admin/search [get]
func (h *Handler) AdminSearch(c *gin.Context) {
	p, ok := platformauth.RequireSuperAdmin(c)
	if !ok {
		return
	}
	data, err := h.svc.AdminSearch(c.Request.Context(), p, "")
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// AdminSearchByValue godoc
// @Summary Search shared applications (super-admin)
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Param value path string true "Filter"
// @Success 200 {object} response.Envelope
// @Router /private/applications/admin/search/{value} [get]
func (h *Handler) AdminSearchByValue(c *gin.Context) {
	p, ok := platformauth.RequireSuperAdmin(c)
	if !ok {
		return
	}
	data, err := h.svc.AdminSearch(c.Request.Context(), p, c.Param("value"))
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// TurnIntoCommon godoc
// @Summary Mark application as shared catalog entry
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Success 200 {object} response.Envelope
// @Router /private/applications/admin/common/{id} [get]
func (h *Handler) TurnIntoCommon(c *gin.Context) {
	p, ok := platformauth.RequireSuperAdmin(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.TurnIntoCommon(c.Request.Context(), p, id); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}
