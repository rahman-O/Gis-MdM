package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/modules/settings/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/settings/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves settings HTTP API.
type Handler struct {
	svc *application.Service
}

func NewHandler(svc *application.Service) *Handler {
	return &Handler{svc: svc}
}

func Register(g *gin.RouterGroup, h *Handler) {
	g.GET("", h.Get)
	g.POST("/misc", h.SaveMisc)
	g.POST("/lang", h.SaveLang)
	g.POST("/design", h.SaveDesign)
	g.GET("/userRole/:roleId", h.GetUserRole)
	g.POST("/userRoles/common", h.SaveUserRolesCommon)
}

func (h *Handler) Get(c *gin.Context) {
	principal, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.ErrorEnvelope(c, "error.permission.denied")
		return
	}
	settings, err := h.svc.Get(c.Request.Context(), principal.CustomerID)
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, settings)
}

func (h *Handler) SaveMisc(c *gin.Context) {
	h.saveBody(c, func(ctx *gin.Context, cid int, body map[string]interface{}) error {
		return h.svc.MergeAndSaveMisc(ctx, cid, body)
	})
}

func (h *Handler) SaveLang(c *gin.Context) {
	h.saveBody(c, func(ctx *gin.Context, cid int, body map[string]interface{}) error {
		return h.svc.MergeAndSaveLang(ctx, cid, body)
	})
}

func (h *Handler) SaveDesign(c *gin.Context) {
	h.saveBody(c, func(ctx *gin.Context, cid int, body map[string]interface{}) error {
		return h.svc.MergeAndSaveDesign(ctx, cid, body)
	})
}

func (h *Handler) saveBody(c *gin.Context, fn func(*gin.Context, int, map[string]interface{}) error) {
	principal, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := fn(c, principal.CustomerID, body); err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

func (h *Handler) GetUserRole(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("roleId"))
	response.OK(c, domain.UserRoleSettings{RoleID: roleID})
}

func (h *Handler) SaveUserRolesCommon(c *gin.Context) {
	response.OK(c, nil)
}
