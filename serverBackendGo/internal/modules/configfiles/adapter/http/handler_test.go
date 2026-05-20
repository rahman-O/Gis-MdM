package http

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type filesDirStub struct{}

func (filesDirStub) FilesDir(*gin.Context, int) (string, error) {
	return "test-files", nil
}

func TestUpload_ok(t *testing.T) {
	tmp := t.TempDir()
	h := NewHandler(tmp, "http://localhost:8080", filesDirStub{})
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		platformauth.WithPrincipal(c, &platformauth.Principal{ID: 1, CustomerID: 1, AuthLoaded: true})
		c.Next()
	})
	h.Register(r.Group("/rest/private/config-files"))

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	fw, _ := w.CreateFormFile("file", "test.txt")
	_, _ = fw.Write([]byte("hello"))
	_ = w.Close()

	req := httptest.NewRequest(http.MethodPost, "/rest/private/config-files", body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
	if _, err := os.Stat(filepath.Join(tmp, "test-files", "test.txt")); err != nil {
		t.Fatalf("file not written: %v", err)
	}
}
