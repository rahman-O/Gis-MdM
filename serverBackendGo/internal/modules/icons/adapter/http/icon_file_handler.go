package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	iconapp "github.com/gis-mdm/server-backend-go/internal/modules/icons/application"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// IconFileHandler serves POST /rest/private/icon-files.
type IconFileHandler struct {
	store    iconapp.IconFileStore
	customer iconapp.CustomerFilesDir
	files    iconapp.UploadedFileInserter
}

func NewIconFileHandler(store iconapp.IconFileStore, customer iconapp.CustomerFilesDir, files iconapp.UploadedFileInserter) *IconFileHandler {
	return &IconFileHandler{store: store, customer: customer, files: files}
}

func (h *IconFileHandler) Register(g *gin.RouterGroup) {
	g.POST("", h.Upload)
}

// Upload godoc
// @Summary Upload icon image file
// @Tags Icons
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Icon image"
// @Success 200 {object} response.Envelope
// @Router /private/icon-files [post]
func (h *IconFileHandler) Upload(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	fh, err := c.FormFile("file")
	if err != nil {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	f, err := fh.Open()
	if err != nil {
		response.ErrorEnvelope(c, "error.invalid.request")
		return
	}
	defer f.Close()
	out, err := iconapp.UploadIconFile(c.Request.Context(), p, f, h.store, h.customer, h.files)
	if err != nil {
		if errors.Is(err, iconapp.ErrIconDimensionInvalid) {
			response.ErrorEnvelope(c, err.Error())
			return
		}
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	response.OK(c, out)
}
