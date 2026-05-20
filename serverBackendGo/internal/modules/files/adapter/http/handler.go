package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	fileapp "github.com/gis-mdm/server-backend-go/internal/modules/files/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/files/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// Handler serves /private/web-ui-files.
type Handler struct {
	svc *fileapp.Service
}

func NewHandler(svc *fileapp.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *gin.RouterGroup) {
	g.GET("/search", h.Search)
	g.GET("/search/:value", h.SearchByValue)
	g.POST("/remove", h.Remove)
	g.POST("/update", h.Update)
	g.POST("", h.Upload)
	g.POST("/raw", h.UploadRaw)
	g.GET("/limit", h.Limit)
	g.GET("/apps/*url", h.AppsByURL)
	g.GET("/configurations/:id", h.FileConfigurations)
	g.POST("/configurations", h.UpdateFileConfigurations)
}

func principal(c *gin.Context) (*platformauth.Principal, bool) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok || p == nil {
		c.Status(403)
		return nil, false
	}
	return p, true
}

func mapErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, fileapp.ErrPermissionDenied):
		response.PermissionDenied(c)
	case errors.Is(err, fileapp.ErrFileUsed):
		response.ErrorEnvelope(c, "error.used.file")
	case errors.Is(err, fileapp.ErrFileExists):
		response.ErrorEnvelope(c, "error.duplicate.file")
	case errors.Is(err, fileapp.ErrSizeLimit):
		response.ErrorEnvelope(c, "error.size.limit.exceeded")
	case errors.Is(err, fileapp.ErrSaveFile):
		response.ErrorEnvelope(c, "error.file.save")
	default:
		response.ErrorEnvelope(c, "error.internal.server")
	}
}

// Search godoc
// @Summary List uploaded files
// @Tags Files
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/web-ui-files/search [get]
func (h *Handler) Search(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	data, err := h.svc.Search(c.Request.Context(), p, "")
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// SearchByValue godoc
// @Summary Search uploaded files
// @Tags Files
// @Produce json
// @Security BearerAuth
// @Param value path string true "filter"
// @Success 200 {object} response.Envelope
// @Router /private/web-ui-files/search/{value} [get]
func (h *Handler) SearchByValue(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	data, err := h.svc.Search(c.Request.Context(), p, c.Param("value"))
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// Remove godoc
// @Summary Remove uploaded file
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/web-ui-files/remove [post]
func (h *Handler) Remove(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var body domain.FileView
	if err := c.ShouldBindJSON(&body); err != nil || body.ID == nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	if err := h.svc.Remove(c.Request.Context(), p, *body.ID, body.FilePath, body.External); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}

// Update godoc
// @Summary Create or update uploaded file metadata
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/web-ui-files/update [post]
func (h *Handler) Update(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var body domain.UploadedFile
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	out, err := h.svc.Commit(c.Request.Context(), p, body)
	if err != nil {
		mapErr(c, err)
		return
	}
	if out != nil && (body.ID == nil || *body.ID == 0) {
		response.OK(c, out)
		return
	}
	response.OK(c, nil)
}

// Upload godoc
// @Summary Upload file with optional APK parse
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "file"
// @Success 200 {object} response.Envelope
// @Router /private/web-ui-files [post]
func (h *Handler) Upload(c *gin.Context) {
	h.upload(c, true)
}

// UploadRaw godoc
// @Summary Upload raw file without APK parse
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "file"
// @Success 200 {object} response.Envelope
// @Router /private/web-ui-files/raw [post]
func (h *Handler) UploadRaw(c *gin.Context) {
	h.upload(c, false)
}

func (h *Handler) upload(c *gin.Context, parseAPK bool) {
	p, ok := principal(c)
	if !ok {
		return
	}
	fh, err := c.FormFile("file")
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	src, err := fh.Open()
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	defer src.Close()
	result, err := h.svc.Upload(c.Request.Context(), p, fh.Filename, src, parseAPK)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, result)
}

// Limit godoc
// @Summary Storage limit for tenant
// @Tags Files
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/web-ui-files/limit [get]
func (h *Handler) Limit(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	data, err := h.svc.GetLimit(c.Request.Context(), p)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// AppsByURL godoc
// @Summary Applications using file URL
// @Tags Files
// @Produce json
// @Security BearerAuth
// @Param url path string true "file URL"
// @Success 200 {object} response.Envelope
// @Router /private/web-ui-files/apps/{url} [get]
func (h *Handler) AppsByURL(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	raw := c.Param("url")
	if len(raw) > 0 && raw[0] == '/' {
		raw = raw[1:]
	}
	data, err := h.svc.GetApplicationsByURL(c.Request.Context(), p, raw)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// FileConfigurations godoc
// @Summary Configurations linked to file
// @Tags Files
// @Produce json
// @Security BearerAuth
// @Param id path int true "file id"
// @Success 200 {object} response.Envelope
// @Router /private/web-ui-files/configurations/{id} [get]
func (h *Handler) FileConfigurations(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.svc.GetFileConfigurations(c.Request.Context(), p, id)
	if err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, data)
}

// UpdateFileConfigurations godoc
// @Summary Update file-configuration links
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Router /private/web-ui-files/configurations [post]
func (h *Handler) UpdateFileConfigurations(c *gin.Context) {
	p, ok := principal(c)
	if !ok {
		return
	}
	var body domain.LinkConfigurationsToFileRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	if err := h.svc.UpdateFileConfigurations(c.Request.Context(), p, body); err != nil {
		mapErr(c, err)
		return
	}
	response.OK(c, nil)
}

// Avoid unused import
var _ = http.StatusOK
