package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	pluginapp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

type Handler struct {
	svc *pluginapp.Service
}

func NewHandler(svc *pluginapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.POST("/private/search", h.Search)
	g.POST("/private/send", h.Send)
	g.DELETE("/private/:id", h.Delete)
	g.GET("/private/purge/:days", h.Purge)
	g.POST("/private/searchTasks", h.SearchTasks)
	g.PUT("/private/task", h.SaveTask)
	g.DELETE("/private/task/:id", h.DeleteTask)
}

func (h *Handler) Search(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var f domain.PushMessageFilter
	_ = c.ShouldBindJSON(&f)
	out, err := h.svc.Search(c.Request.Context(), p, f)
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, map[string]any{"items": out.Items, "totalItemsCount": out.Total})
}

func (h *Handler) Send(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var req domain.PushSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "error.params.missing")
		return
	}
	if err := h.svc.Send(c.Request.Context(), p, req); err != nil {
		if errors.Is(err, pluginapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) Delete(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), p, id); err != nil {
		if errors.Is(err, pluginapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

func (h *Handler) Purge(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	days, _ := strconv.Atoi(c.Param("days"))
	n, err := h.svc.Purge(c.Request.Context(), p, days)
	if err != nil {
		if errors.Is(err, pluginapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, n)
}

func (h *Handler) SearchTasks(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var f domain.PushScheduleFilter
	_ = c.ShouldBindJSON(&f)
	out, err := h.svc.SearchTasks(c.Request.Context(), p, f)
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, out)
}

func (h *Handler) SaveTask(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var task domain.PluginPushSchedule
	if err := c.ShouldBindJSON(&task); err != nil {
		response.ErrorEnvelope(c, "error.params.missing")
		return
	}
	if err := h.svc.SaveTask(c.Request.Context(), p, task); err != nil {
		if errors.Is(err, pluginapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.DeleteTask(c.Request.Context(), p, id); err != nil {
		if errors.Is(err, pluginapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}
