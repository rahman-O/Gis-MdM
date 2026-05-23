package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/application"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/middleware"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves /private/twofactor/* endpoints.
type Handler struct {
	svc *application.TwoFactorService
}

// NewHandler creates the handler.
func NewHandler(svc *application.TwoFactorService) *Handler {
	return &Handler{svc: svc}
}

// Register mounts routes on /twofactor.
func Register(g *gin.RouterGroup, h *Handler) {
	g.GET("/qr/:userId", h.QR)
	g.GET("/verify/:userId/:code", h.Verify)
	g.GET("/set", h.Set)
	g.GET("/reset", h.Reset)
}

// QR godoc
// @Summary Two-factor — QR code PNG
// @Tags TwoFactor
// @Produce png
// @Security BearerAuth
// @Success 200 {file} binary
// @Router /private/twofactor/qr/{userId} [get]
func (h *Handler) QR(c *gin.Context) {
	principal, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		c.Status(http.StatusForbidden)
		return
	}
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil || userID != principal.ID {
		c.Status(http.StatusForbidden)
		return
	}
	png, err := h.svc.QRCodePNG(c.Request.Context(), userID, principal.Login)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK, "image/png", png)
}

// Verify godoc
// @Summary Two-factor — verify TOTP code
// @Tags TwoFactor
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/twofactor/verify/{userId}/{code} [get]
func (h *Handler) Verify(c *gin.Context) {
	principal, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil || userID != principal.ID {
		response.PermissionDenied(c)
		return
	}
	code := c.Param("code")
	if err := h.svc.Verify(c.Request.Context(), userID, code); err != nil {
		if errors.Is(err, application.ErrInvalidTOTP) {
			response.PermissionDenied(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	middleware.ClearTwoFactorPending(c)
	response.OK(c, nil)
}

// Set godoc
// @Summary Two-factor — mark accepted
// @Tags TwoFactor
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/twofactor/set [get]
func (h *Handler) Set(c *gin.Context) {
	principal, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	if err := h.svc.SetAccepted(c.Request.Context(), principal.ID); err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

// Reset godoc
// @Summary Two-factor — reset secret
// @Tags TwoFactor
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/twofactor/reset [get]
func (h *Handler) Reset(c *gin.Context) {
	principal, ok := platformauth.PrincipalFromContext(c)
	if !ok {
		response.PermissionDenied(c)
		return
	}
	if err := h.svc.Reset(c.Request.Context(), principal.ID); err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}
