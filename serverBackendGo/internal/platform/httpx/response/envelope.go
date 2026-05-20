package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperr "github.com/gis-mdm/server-backend-go/internal/shared/errors"
)

// Envelope mirrors the legacy Headwind MDM JSON response shape.
type Envelope struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{Status: "OK", Data: data})
}

func ErrorEnvelope(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Envelope{Status: "ERROR", Message: message})
}

func DuplicateEntity(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Envelope{Status: "ERROR", Message: message})
}

func ObjectNotFound(c *gin.Context) {
	c.JSON(http.StatusOK, Envelope{Status: "ERROR", Message: "error.notfound.object"})
}

func PermissionDenied(c *gin.Context) {
	c.JSON(http.StatusOK, Envelope{Status: "ERROR", Message: "error.permission.denied"})
}

func Error(c *gin.Context, err *apperr.AppError) {
	status, body := mapError(err)
	c.JSON(status, body)
}

func NotImplemented(c *gin.Context, message string) {
	Error(c, apperr.NotImplemented(message))
}

func mapError(err *apperr.AppError) (int, Envelope) {
	switch err.Code {
	case apperr.CodeNotFound:
		return http.StatusNotFound, Envelope{Status: "ERROR", Message: err.Message}
	case apperr.CodeUnauthorized:
		return http.StatusUnauthorized, Envelope{Status: "ERROR", Message: err.Message}
	case apperr.CodeForbidden:
		return http.StatusForbidden, Envelope{Status: "ERROR", Message: err.Message}
	case apperr.CodeValidation:
		return http.StatusBadRequest, Envelope{Status: "ERROR", Message: err.Message}
	case apperr.CodeNotImplemented:
		return http.StatusNotImplemented, Envelope{Status: "ERROR", Message: err.Message}
	default:
		return http.StatusInternalServerError, Envelope{Status: "ERROR", Message: err.Message}
	}
}
