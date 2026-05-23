package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	dlapp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/devicelog/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/devicelog/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

type Handler struct {
	svc *dlapp.Service
}

func NewHandler(svc *dlapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(plugins *gin.RouterGroup, engine *gin.Engine) {
	settings := plugins.Group("/devicelog-plugin-settings")
	settings.GET("/private", h.GetSettings)
	settings.PUT("/private", h.PutSettings)
	settings.PUT("/private/rule", h.PutRule)
	settings.DELETE("/private/rule/:id", h.DeleteRule)
	log := plugins.Group("/devicelog/log")
	log.POST("/private/search", h.Search)
	engine.POST("/rest/plugins/devicelog/log/list/:deviceNumber", h.Upload)
}

func (h *Handler) GetSettings(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	out, err := h.svc.GetSettings(c.Request.Context(), p)
	if err != nil {
		if errors.Is(err, dlapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, out)
}

func (h *Handler) PutSettings(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var body domain.Settings
	_ = c.ShouldBindJSON(&body)
	if err := h.svc.SaveSettings(c.Request.Context(), p, body); err != nil {
		if errors.Is(err, dlapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

func (h *Handler) PutRule(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var body domain.Rule
	_ = c.ShouldBindJSON(&body)
	id, err := h.svc.SaveRule(c.Request.Context(), p, body)
	if err != nil {
		if errors.Is(err, dlapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *Handler) DeleteRule(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.DeleteRule(c.Request.Context(), p, id); err != nil {
		if errors.Is(err, dlapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

func (h *Handler) Search(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var f domain.LogFilter
	_ = c.ShouldBindJSON(&f)
	out, err := h.svc.SearchLogs(c.Request.Context(), p, f)
	if err != nil {
		if errors.Is(err, dlapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, out)
}

func (h *Handler) Upload(c *gin.Context) {
	var rows []domain.UploadRecord
	if err := c.ShouldBindJSON(&rows); err != nil {
		response.ErrorEnvelope(c, "error.params.missing")
		return
	}
	if err := h.svc.UploadLogs(c.Request.Context(), c.Param("deviceNumber"), rows); err != nil {
		if errors.Is(err, dlapp.ErrDeviceNotFound) {
			response.ErrorEnvelope(c, "error.device.notfound")
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}
