package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves signup endpoints.
type Handler struct {
	svc *application.SignupService
}

// NewHandler creates the handler.
func NewHandler(svc *application.SignupService) *Handler {
	return &Handler{svc: svc}
}

// Register mounts routes on /signup.
func Register(g *gin.RouterGroup, h *Handler) {
	g.POST("/verifyEmail", h.VerifyEmail)
	g.GET("/verifyToken/:token", h.VerifyToken)
	g.POST("/complete", h.Complete)
	g.GET("/canSignup", h.CanSignup)
}

func (h *Handler) VerifyEmail(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Language string `json:"language"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	if err := h.svc.VerifyEmail(c.Request.Context(), body.Email, body.Language); err != nil {
		if errors.Is(err, application.ErrDuplicateEmail) {
			response.DuplicateEntity(c, "signup.email.used")
			return
		}
		response.ErrorEnvelope(c, "")
		return
	}
	response.OK(c, nil)
}

func (h *Handler) VerifyToken(c *gin.Context) {
	p, err := h.svc.VerifyToken(c.Request.Context(), c.Param("token"))
	if err != nil {
		if errors.Is(err, application.ErrSignupTokenNotFound) {
			response.ObjectNotFound(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, p)
}

func (h *Handler) Complete(c *gin.Context) {
	var body struct {
		Token       string `json:"token"`
		Name        string `json:"name"`
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		Company     string `json:"company"`
		Description string `json:"description"`
		Passwd      string `json:"passwd"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorEnvelope(c, "")
		return
	}
	err := h.svc.Complete(c.Request.Context(), domain.SignupComplete{
		Token: body.Token, Name: body.Name, FirstName: body.FirstName, LastName: body.LastName,
		Company: body.Company, Description: body.Description, PasswordMD5: body.Passwd,
	})
	if err != nil {
		if errors.Is(err, application.ErrDuplicateCustomer) {
			response.DuplicateEntity(c, "error.duplicate.customer.name")
			return
		}
		if errors.Is(err, application.ErrSignupTokenNotFound) {
			response.ObjectNotFound(c)
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, nil)
}

func (h *Handler) CanSignup(c *gin.Context) {
	if h.svc.Enabled() {
		response.OK(c, nil)
		return
	}
	response.ErrorEnvelope(c, "")
}
