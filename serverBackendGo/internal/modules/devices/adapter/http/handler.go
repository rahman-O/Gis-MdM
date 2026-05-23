package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	devapp "github.com/gis-mdm/server-backend-go/internal/modules/devices/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/devices/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves /private/devices/* endpoints.
type Handler struct {
	svc *devapp.Service
}

func NewHandler(svc *devapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.POST("/search", h.Search)
	g.GET("/number/:number", h.GetByNumber)
	g.PUT("", h.Save)
	g.DELETE("/:id", h.Delete)
	g.POST("/deleteBulk", h.DeleteBulk)
	g.POST("/groupBulk", h.GroupBulk)
	g.POST("/autocomplete", h.Autocomplete)
	g.POST("/:id/description", h.UpdateDescription)
	g.GET("/:id/applicationSettings", h.GetAppSettings)
	g.POST("/:id/applicationSettings", h.SaveAppSettings)
	g.POST("/:id/applicationSettings/notify", h.NotifyAppSettings)
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
	case errors.Is(err, devapp.ErrPermissionDenied):
		response.PermissionDenied(c)
	case errors.Is(err, devapp.ErrDeviceExists):
		response.DuplicateEntity(c, "error.duplicate.device")
	case errors.Is(err, devapp.ErrDeviceNotFound):
		response.ObjectNotFound(c)
	case errors.Is(err, devapp.ErrDeviceLimit):
		response.ErrorEnvelope(c, "error.device.limit")
	default:
		response.ErrorEnvelope(c, "error.internal.server")
	}
}

// Search godoc
// @Summary Search devices
// @Tags Devices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/devices/search [post]
func (h *Handler) Search(c *gin.Context) {
	p, ok := principal(c)
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

// GetByNumber godoc
// @Summary Get device by number
// @Tags Devices
// @Produce json
// @Security BearerAuth
// @Param number path string true "Device number"
// @Success 200 {object} response.Envelope
// @Router /private/devices/number/{number} [get]
func (h *Handler) GetByNumber(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	data, err := h.svc.GetByNumber(c.Request.Context(), p, c.Param("number"))
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// Save godoc
// @Summary Create or update device
// @Tags Devices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/devices [put]
func (h *Handler) Save(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var d domain.SaveDevice
	if err := c.ShouldBindJSON(&d); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.Save(c.Request.Context(), p, d); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete godoc
// @Summary Delete device
// @Tags Devices
// @Produce json
// @Security BearerAuth
// @Param id path int true "Device ID"
// @Success 200 {object} response.Envelope
// @Router /private/devices/{id} [delete]
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

// DeleteBulk godoc
// @Summary Bulk delete devices
// @Tags Devices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/devices/deleteBulk [post]
func (h *Handler) DeleteBulk(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var req domain.BulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.DeleteBulk(c.Request.Context(), p, req); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}

// GroupBulk godoc
// @Summary Bulk update device groups
// @Tags Devices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/devices/groupBulk [post]
func (h *Handler) GroupBulk(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var req domain.GroupBulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.GroupBulk(c.Request.Context(), p, req); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}

// Autocomplete godoc
// @Summary Device number autocomplete
// @Tags Devices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/devices/autocomplete [post]
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

// UpdateDescription godoc
// @Summary Update device description
// @Tags Devices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Device ID"
// @Success 200 {object} response.Envelope
// @Router /private/devices/{id}/description [post]
func (h *Handler) UpdateDescription(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var desc string
	if err := c.ShouldBindJSON(&desc); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.UpdateDescription(c.Request.Context(), p, id, desc); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}

// GetAppSettings godoc
// @Summary Get device application settings
// @Tags Devices
// @Produce json
// @Security BearerAuth
// @Param id path int true "Device ID"
// @Success 200 {object} response.Envelope
// @Router /private/devices/{id}/applicationSettings [get]
func (h *Handler) GetAppSettings(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.svc.GetAppSettings(c.Request.Context(), p, id)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// SaveAppSettings godoc
// @Summary Save device application settings
// @Tags Devices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Device ID"
// @Success 200 {object} response.Envelope
// @Router /private/devices/{id}/applicationSettings [post]
func (h *Handler) SaveAppSettings(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var settings []domain.AppSetting
	if err := c.ShouldBindJSON(&settings); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.SaveAppSettings(c.Request.Context(), p, id, settings); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}

// NotifyAppSettings godoc
// @Summary Notify device of application settings change
// @Tags Devices
// @Produce json
// @Security BearerAuth
// @Param id path int true "Device ID"
// @Success 200 {object} response.Envelope
// @Router /private/devices/{id}/applicationSettings/notify [post]
func (h *Handler) NotifyAppSettings(c *gin.Context) {
	if _, ok := principal(c); !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.NotifyAppSettings(c.Request.Context(), id); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}
