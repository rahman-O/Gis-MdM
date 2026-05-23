package http

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	hintapp "github.com/gis-mdm/server-backend-go/internal/modules/hints/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/hints/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves /private/hints/* endpoints.
type Handler struct {
	svc *hintapp.Service
}

func NewHandler(svc *hintapp.Service) *Handler {
	return &Handler{svc: svc}
}

// Register mounts routes on /hints.
func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/history", h.GetHistory)
	g.POST("/history", h.MarkShown)
	g.POST("/enable", h.Enable)
	g.POST("/disable", h.Disable)
}

func principal(c *gin.Context) (*platformauth.Principal, bool) {
	return platformauth.PrincipalFromContext(c)
}

// GetHistory godoc
// @Summary Get shown hint keys
// @Tags Hints
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/hints/history [get]
func (h *Handler) GetHistory(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		c.Status(403)
		return
	}
	keys, err := h.svc.GetHistory(c.Request.Context(), p)
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, keys)
}

// MarkShown godoc
// @Summary Mark hint as shown
// @Tags Hints
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/hints/history [post]
func (h *Handler) MarkShown(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		c.Status(403)
		return
	}
	key, err := parseHintKeyBody(c)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyHintKey) {
			response.ErrorEnvelope(c, "error.hint.empty")
			return
		}
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.MarkShown(c.Request.Context(), p, key); err != nil {
		if errors.Is(err, domain.ErrEmptyHintKey) {
			response.ErrorEnvelope(c, "error.hint.empty")
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

// Enable godoc
// @Summary Enable hints (clear history)
// @Tags Hints
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/hints/enable [post]
func (h *Handler) Enable(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		c.Status(403)
		return
	}
	if err := h.svc.Enable(c.Request.Context(), p); err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

// Disable godoc
// @Summary Disable hints (mark all catalog keys shown)
// @Tags Hints
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/hints/disable [post]
func (h *Handler) Disable(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		c.Status(403)
		return
	}
	if err := h.svc.Disable(c.Request.Context(), p); err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

func parseHintKeyBody(c *gin.Context) (string, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return "", err
	}
	raw := strings.TrimSpace(string(body))
	if raw == "" {
		return "", domain.ErrEmptyHintKey
	}
	// JSON string: "hint.step.1"
	var asString string
	if err := json.Unmarshal(body, &asString); err == nil && asString != "" {
		_, err := domain.ValidateHintKey(asString)
		return asString, err
	}
	// JSON object: {"hintKey":"..."}
	var wrapped struct {
		HintKey string `json:"hintKey"`
	}
	if err := json.Unmarshal(body, &wrapped); err == nil && wrapped.HintKey != "" {
		_, err := domain.ValidateHintKey(wrapped.HintKey)
		return wrapped.HintKey, err
	}
	_, err = domain.ValidateHintKey(raw)
	return raw, err
}
