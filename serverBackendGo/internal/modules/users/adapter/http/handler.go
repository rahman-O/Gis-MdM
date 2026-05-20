package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/application"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves /private/users/* endpoints.
type Handler struct {
	svc *application.UsersService
}

// NewHandler creates the handler.
func NewHandler(svc *application.UsersService) *Handler {
	return &Handler{svc: svc}
}

// Register mounts routes on /users.
func Register(g *gin.RouterGroup, h *Handler) {
	g.GET("/current", h.Current)
}

func (h *Handler) Current(c *gin.Context) {
	principal, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.OK(c, nil)
		return
	}
	view, err := h.svc.CurrentUser(c.Request.Context(), principal.ID)
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, view)
}
