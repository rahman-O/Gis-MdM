package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/application"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves password reset endpoints.
type Handler struct {
	svc *application.PasswordResetService
}

// NewHandler creates the handler.
func NewHandler(svc *application.PasswordResetService) *Handler {
	return &Handler{svc: svc}
}

// Register mounts routes on /passwordReset.
func Register(g *gin.RouterGroup, h *Handler) {
	g.GET("/settings/:token", h.GetSettings)
	g.POST("/reset", h.Reset)
	g.GET("/recover/:username", h.Recover)
	g.GET("/canRecover", h.CanRecover)
}

// GetSettings godoc
// @Summary Password reset — settings for token
// @Tags PasswordReset
// @Produce json
// @Success 200 {object} response.Envelope
// @Router /public/passwordReset/settings/{token} [get]
func (h *Handler) GetSettings(c *gin.Context) {
	token := c.Param("token")
	settings, err := h.svc.ResetSettingsByToken(c.Request.Context(), token)
	if err != nil {
		if errors.Is(err, application.ErrUserNotFound) {
			response.ErrorEnvelope(c, "error.user.not.found")
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, settings)
}

// Reset godoc
// @Summary Password reset — set new password
// @Tags PasswordReset
// @Accept json
// @Produce json
// @Success 200 {object} response.Envelope
// @Router /public/passwordReset/reset [post]
func (h *Handler) Reset(c *gin.Context) {
	var body struct {
		PasswordResetToken string `json:"passwordResetToken"`
		NewPassword        string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	view, err := h.svc.ResetPassword(c.Request.Context(), body.PasswordResetToken, body.NewPassword)
	if err != nil {
		if errors.Is(err, application.ErrUserNotFound) {
			response.ErrorEnvelope(c, "error.user.not.found")
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, view)
}

// Recover godoc
// @Summary Password reset — send recovery email
// @Tags PasswordReset
// @Produce json
// @Success 200 {object} response.Envelope
// @Router /public/passwordReset/recover/{username} [get]
func (h *Handler) Recover(c *gin.Context) {
	username := c.Param("username")
	if err := h.svc.Recover(c.Request.Context(), username); err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

// CanRecover godoc
// @Summary Password reset — can recover (deprecated)
// @Tags PasswordReset
// @Produce json
// @Success 200 {object} response.Envelope
// @Router /public/passwordReset/canRecover [get]
func (h *Handler) CanRecover(c *gin.Context) {
	response.OK(c, nil)
}
