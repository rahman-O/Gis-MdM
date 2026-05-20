package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	diapp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/deviceinfo/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/deviceinfo/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

type Handler struct {
	svc *diapp.Service
}

func NewHandler(svc *diapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(plugins *gin.RouterGroup, engine *gin.Engine) {
	settings := plugins.Group("/deviceinfo-plugin-settings")
	settings.GET("/private", h.GetSettings)
	settings.PUT("/private", h.PutSettings)
	info := plugins.Group("/deviceinfo")
	info.GET("/private/:deviceNumber", h.GetDetail)
	info.POST("/private/search/dynamic", h.SearchDynamic)
	engine.PUT("/rest/plugins/deviceinfo/deviceinfo/public/:deviceNumber", h.PutPublic)
}

func (h *Handler) GetSettings(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	out, err := h.svc.GetSettings(c.Request.Context(), p)
	if err != nil {
		if errors.Is(err, diapp.ErrPermissionDenied) {
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
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorEnvelope(c, "error.params.missing")
		return
	}
	if err := h.svc.SaveSettings(c.Request.Context(), p, body); err != nil {
		if errors.Is(err, diapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

func (h *Handler) PutPublic(c *gin.Context) {
	var items []domain.DynamicInfo
	if err := c.ShouldBindJSON(&items); err != nil {
		response.ErrorEnvelope(c, "error.params.missing")
		return
	}
	if err := h.svc.SavePublicDynamic(c.Request.Context(), c.Param("deviceNumber"), items); err != nil {
		if errors.Is(err, diapp.ErrDeviceNotFound) {
			response.ErrorEnvelope(c, "error.device.notfound")
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

func (h *Handler) GetDetail(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	out, err := h.svc.GetDeviceDetail(c.Request.Context(), p, c.Param("deviceNumber"))
	if err != nil {
		if errors.Is(err, diapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		if errors.Is(err, diapp.ErrDeviceNotFound) {
			response.ErrorEnvelope(c, "error.device.notfound")
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, out)
}

func (h *Handler) SearchDynamic(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var f domain.DynamicSearchFilter
	_ = c.ShouldBindJSON(&f)
	out, err := h.svc.SearchDynamic(c.Request.Context(), p, f)
	if err != nil {
		if errors.Is(err, diapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, out)
}
