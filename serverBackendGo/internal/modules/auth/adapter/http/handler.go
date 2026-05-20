package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/application"
	"github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/middleware"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
	apperr "github.com/gis-mdm/server-backend-go/internal/shared/errors"
)

// Handler serves auth HTTP endpoints.
type Handler struct {
	svc *application.Service
}

// NewHandler creates an auth HTTP handler.
func NewHandler(svc *application.Service) *Handler {
	return &Handler{svc: svc}
}

// Options godoc
// @Summary Login options
// @Tags Authentication
// @Produce json
// @Success 200 {object} response.Envelope
// @Router /public/auth/options [get]
func (h *Handler) Options(c *gin.Context) {
	opts := h.svc.Options(c.Request.Context())
	response.OK(c, opts)
}

// Login godoc
// @Summary Session login
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Credentials (password: raw or MD5 uppercase hex)"
// @Success 200 {object} response.Envelope
// @Router /public/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	view, twoFactorPending, err := h.svc.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		var fail application.AuthFailure
		if errors.As(err, &fail) {
			response.ErrorEnvelope(c, "")
			return
		}
		response.Error(c, apperr.Internal("login failed", err))
		return
	}
	middleware.SessionStore(c, &auth.Principal{
		ID: view.ID, Login: view.Login, AuthToken: view.AuthToken, CustomerID: view.CustomerID,
		PasswordReset: view.PasswordReset,
	}, twoFactorPending)
	response.OK(c, view)
}

// Logout godoc
// @Summary Logout
// @Tags Authentication
// @Success 204
// @Router /public/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	h.svc.Logout()
	middleware.SessionClear(c)
	c.Status(http.StatusNoContent)
}

// JWTLogin godoc
// @Summary JWT login
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Credentials (password: raw or MD5 uppercase hex)"
// @Success 200 {object} JWTResultDTO
// @Failure 400
// @Failure 401
// @Header 200 {string} Authorization "Bearer token"
// @Router /public/jwt/login [post]
func (h *Handler) JWTLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	result, authHeader, err := h.svc.JWTLogin(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		var bad application.BadRequest
		if errors.As(err, &bad) {
			c.Status(http.StatusBadRequest)
			return
		}
		var unauth application.Unauthorized
		if errors.As(err, &unauth) {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}
	if authHeader != "" {
		c.Header("Authorization", authHeader)
	}
	c.JSON(http.StatusOK, result)
}
