package http

import (
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
	"github.com/gis-mdm/server-backend-go/internal/platform/storage"
)

// Handler serves POST /private/config-files.
type Handler struct {
	filesDir string
	baseURL  string
	db       customerFilesDir
}

type customerFilesDir interface {
	FilesDir(ctx *gin.Context, customerID int) (string, error)
}

// NewHandler creates the config file upload handler.
func NewHandler(filesDir, baseURL string, db customerFilesDir) *Handler {
	return &Handler{filesDir: filesDir, baseURL: baseURL, db: db}
}

// Register mounts routes on /config-files.
func (h *Handler) Register(g *gin.RouterGroup) {
	g.POST("", h.Upload)
}

// UploadResult mirrors Java UploadedFile / FileUploadResult subset for React.
type UploadResult struct {
	CustomerID int    `json:"customerId,omitempty"`
	FilePath   string `json:"filePath"`
	URL        string `json:"url,omitempty"`
	Name       string `json:"name,omitempty"`
}

// Upload godoc
// @Summary Upload configuration file asset
// @Tags ConfigFiles
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File"
// @Success 200 {object} response.Envelope
// @Router /private/config-files [post]
func (h *Handler) Upload(c *gin.Context) {
	p, ok := platformauth.PrincipalFromContext(c)
	if !ok || p == nil {
		c.Status(403)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	filesDir, err := h.db.FilesDir(c, p.CustomerID)
	if err != nil || filesDir == "" {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	targetDir := filepath.Join(h.filesDir, filesDir)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	safeName := filepath.Base(file.Filename)
	dest := filepath.Join(targetDir, safeName)
	src, err := file.Open()
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	defer src.Close()
	out, err := os.Create(dest)
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	if _, err = io.Copy(out, src); err != nil {
		out.Close()
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	out.Close()
	url := storage.BuildPublicURL(h.baseURL, filesDir, safeName)
	result := UploadResult{
		CustomerID: p.CustomerID,
		FilePath:   safeName,
		URL:        url,
		Name:       safeName,
	}
	response.OK(c, result)
}
