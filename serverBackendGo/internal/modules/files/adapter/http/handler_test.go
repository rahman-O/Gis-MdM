package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	fileapp "github.com/gis-mdm/server-backend-go/internal/modules/files/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/files/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/platform/storage"
)

type httpFileStub struct{}

func (httpFileStub) List(context.Context, int, string) ([]domain.UploadedFile, error) {
	return []domain.UploadedFile{{ID: intPtr(1)}}, nil
}
func (httpFileStub) GetByID(context.Context, int, int) (*domain.UploadedFile, error) { return nil, nil }
func (httpFileStub) Insert(context.Context, *domain.UploadedFile) error               { return nil }
func (httpFileStub) Update(context.Context, *domain.UploadedFile) error               { return nil }
func (httpFileStub) Delete(context.Context, int) error                                { return nil }
func (httpFileStub) IsUsedByConfiguration(context.Context, int) (bool, error)         { return false, nil }
func (httpFileStub) IsUsedByIcon(context.Context, int) (bool, error)                { return false, nil }
func (httpFileStub) UsingConfigurationNames(context.Context, int, int) ([]string, error) {
	return nil, nil
}
func (httpFileStub) UsingIconNames(context.Context, int, int) ([]string, error) { return nil, nil }
func (httpFileStub) GetFileConfigurations(context.Context, int, int, int) ([]domain.FileConfigurationLink, error) {
	return nil, nil
}
func (httpFileStub) DeleteConfigurationFile(context.Context, int) error { return nil }
func (httpFileStub) InsertConfigurationFile(context.Context, int, int, string) error {
	return nil
}
func (httpFileStub) CountByPath(context.Context, int, *int, string) (int64, error) { return 0, nil }

type httpCustStub struct{}

func (httpCustStub) GetMeta(context.Context, int) (*domain.CustomerMeta, error) {
	return &domain.CustomerMeta{}, nil
}
func (httpCustStub) CountCustomers(context.Context) (int, error) { return 1, nil }

func intPtr(n int) *int { return &n }

func setupFilesRouter(h *Handler, perms ...string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	g := r.Group("/private/web-ui-files")
	g.Use(func(c *gin.Context) {
		platformauth.WithPrincipal(c, &platformauth.Principal{
			CustomerID: 1, AuthLoaded: true, Permissions: perms,
		})
		c.Next()
	})
	h.Register(g)
	return r
}

func TestSearchOKWithFilesPermission(t *testing.T) {
	store := storage.NewLocalStore(t.TempDir())
	svc := fileapp.NewService(httpFileStub{}, httpCustStub{}, nil, store, "http://localhost:8080", nil)
	r := setupFilesRouter(NewHandler(svc), "files")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/private/web-ui-files/search", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"status":"OK"`) {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestSearchDeniedWithoutPermission(t *testing.T) {
	svc := fileapp.NewService(httpFileStub{}, httpCustStub{}, nil, nil, "http://localhost:8080", nil)
	r := setupFilesRouter(NewHandler(svc), "settings")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/private/web-ui-files/search", nil))
	if !strings.Contains(w.Body.String(), "error.permission.denied") {
		t.Fatalf("body %s", w.Body.String())
	}
}
