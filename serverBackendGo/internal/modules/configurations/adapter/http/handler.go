package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	cfgapp "github.com/gis-mdm/server-backend-go/internal/modules/configurations/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves /private/configurations/* endpoints.
type Handler struct {
	svc *cfgapp.Service
}

func NewHandler(svc *cfgapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/search", h.Search)
	g.GET("/search/:value", h.SearchByValue)
	g.GET("/list", h.List)
	g.POST("/autocomplete", h.Autocomplete)
	g.PUT("", h.Save)
	g.PUT("/copy", h.Copy)
	g.DELETE("/:id", h.Delete)
	g.GET("/applications", h.ListAllApplications)
	g.GET("/applications/:id", h.ListConfigurationApplications)
	g.PUT("/application/upgrade", h.UpgradeApplication)
	g.GET("/:id", h.GetByID)
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
	case errors.Is(err, cfgapp.ErrPermissionDenied):
		response.PermissionDenied(c)
	case errors.Is(err, cfgapp.ErrDuplicateConfiguration):
		response.DuplicateEntity(c, "error.duplicate.configuration")
	case errors.Is(err, cfgapp.ErrNotEmptyConfiguration):
		response.ErrorEnvelope(c, "error.notempty.configuration")
	case errors.Is(err, cfgapp.ErrConfigurationNotFound):
		response.ErrorEnvelope(c, "error.notfound.configuration")
	default:
		response.ErrorEnvelope(c, "error.internal.server")
	}
}

func okConfigAppList(c *gin.Context, data []domain.ConfigurationApplication) {
	if data == nil {
		data = []domain.ConfigurationApplication{}
	}
	response.OK(c, data)
}

// List godoc
// @Summary List configuration names for dropdowns
// @Tags Configurations
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/configurations/list [get]
func (h *Handler) List(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	data, err := h.svc.ListNames(c.Request.Context(), p)
	if err != nil {
		mapErr(c, err)
		return
	}
	if data == nil {
		data = []domain.LookupItem{}
	}
	response.OK(c, data)
}

// Search godoc
// @Summary List configurations
// @Tags Configurations
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/configurations/search [get]
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
// @Summary Search configurations by name
// @Tags Configurations
// @Produce json
// @Security BearerAuth
// @Param value path string true "Filter"
// @Success 200 {object} response.Envelope
// @Router /private/configurations/search/{value} [get]
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
// @Summary Configuration autocomplete
// @Tags Configurations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/configurations/autocomplete [post]
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
// @Summary Get configuration for editor
// @Tags Configurations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Configuration ID"
// @Success 200 {object} response.Envelope
// @Router /private/configurations/{id} [get]
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
	if data == nil {
		response.ErrorEnvelope(c, "error.notfound.configuration")
		return
	}
	response.OK(c, domain.ConfigurationResponseMap(data))
}

// Save godoc
// @Summary Create or update configuration
// @Tags Configurations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/configurations [put]
func (h *Handler) Save(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	raw, err := c.GetRawData()
	if err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	cfg, err := domain.ParseConfigurationBody(raw)
	if err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	data, err := h.svc.Save(c.Request.Context(), p, cfg)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, domain.ConfigurationResponseMap(data))
}

// Delete godoc
// @Summary Delete configuration
// @Tags Configurations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Configuration ID"
// @Success 200 {object} response.Envelope
// @Router /private/configurations/{id} [delete]
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

// Copy godoc
// @Summary Copy configuration
// @Tags Configurations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/configurations/copy [put]
func (h *Handler) Copy(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var req domain.CopyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	id, err := h.svc.Copy(c.Request.Context(), p, req)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, gin.H{"id": id})
}

// ListAllApplications godoc
// @Summary Applications available for configuration editor
// @Tags Configurations
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/configurations/applications [get]
func (h *Handler) ListAllApplications(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	data, err := h.svc.ListAllApplicationsForPicker(c.Request.Context(), p)
	if err != nil {
		mapErr(c, err)
		return
	}
	okConfigAppList(c, data)
}

// ListConfigurationApplications godoc
// @Summary Applications linked to configuration
// @Tags Configurations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Configuration ID"
// @Success 200 {object} response.Envelope
// @Router /private/configurations/applications/{id} [get]
func (h *Handler) ListConfigurationApplications(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.svc.ListConfigurationApplications(c.Request.Context(), p, id)
	if err != nil {
		mapErr(c, err)
		return
	}
	okConfigAppList(c, data)
}

// UpgradeApplication godoc
// @Summary Upgrade application on configuration to latest version
// @Tags Configurations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/configurations/application/upgrade [put]
func (h *Handler) UpgradeApplication(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var req domain.UpgradeApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.UpgradeApplication(c.Request.Context(), p, req); err != nil {
		mapErr(c, err)
		return
	}
	cfg, _ := h.svc.GetByID(c.Request.Context(), p, req.ConfigurationID)
	if cfg == nil {
		response.OK(c, nil)
		return
	}
	response.OK(c, domain.ConfigurationResponseMap(cfg))
}
