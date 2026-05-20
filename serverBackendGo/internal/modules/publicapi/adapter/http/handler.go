package http

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	pubapp "github.com/gis-mdm/server-backend-go/internal/modules/publicapi/application"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

type Handler struct {
	svc *pubapp.Service
}

func NewHandler(svc *pubapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/name", h.Name)
	g.GET("/logo", h.Logo)
	g.POST("/applications/upload", h.UploadApplication)
}

// Name godoc
// @Summary Rebranding metadata
// @Tags PublicAPI
// @Produce json
// @Success 200 {object} response.Envelope
// @Router /public/name [get]
func (h *Handler) Name(c *gin.Context) {
	response.OK(c, h.svc.GetName())
}

// Logo godoc
// @Summary Rebranding logo
// @Tags PublicAPI
// @Produce png
// @Success 200 {file} binary
// @Router /public/logo [get]
func (h *Handler) Logo(c *gin.Context) {
	path := h.svc.LogoPath()
	if path != "" {
		if _, err := os.Stat(path); err == nil {
			c.Header("Cache-Control", "no-cache")
			c.File(path)
			return
		}
	}
	c.Redirect(http.StatusFound, "/images/logo.png")
}

// UploadApplication godoc
// @Summary AppList application upload
// @Tags PublicAPI
// @Accept multipart/form-data
// @Produce json
// @Param file formData file false "APK file"
// @Param app formData string true "UploadAppRequest JSON"
// @Success 200 {object} response.Envelope
// @Router /public/applications/upload [post]
func (h *Handler) UploadApplication(c *gin.Context) {
	appJSON := c.PostForm("app")
	var r io.Reader
	var fileName string
	fh, err := c.FormFile("file")
	if err == nil && fh != nil {
		fileName = fh.Filename
		src, err := fh.Open()
		if err != nil {
			response.ErrorEnvelope(c, "error.internal.server")
			return
		}
		defer src.Close()
		r = src
	}
	if err := h.svc.UploadApplication(c.Request.Context(), appJSON, fileName, r); err != nil {
		switch {
		case errors.Is(err, pubapp.ErrInvalidHash):
			response.ErrorEnvelope(c, "Invalid hash")
		case errors.Is(err, pubapp.ErrDeviceNotFound):
			response.ErrorEnvelope(c, "error.notfound.device")
		case errors.Is(err, pubapp.ErrDuplicateApp):
			response.ErrorEnvelope(c, "error.duplicate.application")
		case errors.Is(err, pubapp.ErrPermissionDenied):
			response.PermissionDenied(c)
		case errors.Is(err, pubapp.ErrMissingParams):
			response.ErrorEnvelope(c, "error.params.missing")
		default:
			response.ErrorEnvelope(c, "error.internal.server")
		}
		return
	}
	response.OK(c, nil)
}
