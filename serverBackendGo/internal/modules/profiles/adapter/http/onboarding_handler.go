package http

import (
	"github.com/gin-gonic/gin"
	profileapp "github.com/gis-mdm/server-backend-go/internal/modules/profiles/application"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// OnboardingHandler serves GET /private/onboarding/status.
type OnboardingHandler struct {
	svc *profileapp.OnboardingService
}

func NewOnboardingHandler(svc *profileapp.OnboardingService) *OnboardingHandler {
	return &OnboardingHandler{svc: svc}
}

func (h *OnboardingHandler) Register(g *gin.RouterGroup) {
	g.GET("/status", h.Status)
}

func (h *OnboardingHandler) Status(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	data, err := h.svc.Status(c.Request.Context(), p)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}
