package http

import (
	"github.com/gin-gonic/gin"
	cfgpostgres "github.com/gis-mdm/server-backend-go/internal/modules/configurations/adapter/persistence/postgres"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves /private/configurations/list.
type Handler struct {
	repo *cfgpostgres.ConfigRepository
}

func NewHandler(repo *cfgpostgres.ConfigRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/list", h.List)
}

// List godoc
// @Summary List configurations for tenant
// @Tags Configurations
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/configurations/list [get]
func (h *Handler) List(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok || p == nil {
		c.Status(403)
		return
	}
	data, err := h.repo.ListByCustomer(c.Request.Context(), p.CustomerID)
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	if data == nil {
		data = []cfgpostgres.ConfigListItem{}
	}
	response.OK(c, data)
}
