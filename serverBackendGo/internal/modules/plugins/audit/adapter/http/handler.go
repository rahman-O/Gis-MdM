package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	auditapp "github.com/gis-mdm/server-backend-go/internal/modules/plugins/audit/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/audit/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

type Handler struct {
	svc *auditapp.Service
}

func NewHandler(svc *auditapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.POST("/private/log/search", h.Search)
}

func (h *Handler) Search(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	var f domain.AuditLogFilter
	_ = c.ShouldBindJSON(&f)
	out, err := h.svc.Search(c.Request.Context(), p, f)
	if err != nil {
		if errors.Is(err, auditapp.ErrPermissionDenied) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, out)
}
