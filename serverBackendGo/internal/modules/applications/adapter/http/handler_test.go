package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	appapp "github.com/gis-mdm/server-backend-go/internal/modules/applications/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/applications/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type httpAppStub struct{}

func (httpAppStub) Search(context.Context, int) ([]domain.Application, error) {
	return []domain.Application{{ID: intPtr(1)}}, nil
}
func (httpAppStub) SearchByValue(context.Context, int, string) ([]domain.Application, error) {
	return nil, nil
}
func (httpAppStub) GetByID(context.Context, int, int) (*domain.Application, error) { return nil, nil }
func (httpAppStub) ListVersions(context.Context, int, int) ([]domain.ApplicationVersion, error) {
	return []domain.ApplicationVersion{{ID: intPtr(10)}}, nil
}
func (httpAppStub) SaveAndroid(context.Context, int, domain.Application) (*domain.Application, error) {
	return nil, nil
}
func (httpAppStub) SaveWeb(context.Context, int, domain.Application) (*domain.Application, error) {
	return nil, nil
}
func (httpAppStub) SaveVersion(context.Context, int, domain.ApplicationVersion) (*domain.ApplicationVersion, error) {
	return nil, nil
}
func (httpAppStub) DeleteApp(context.Context, int, int) error { return nil }
func (httpAppStub) DeleteVersion(context.Context, int, int) error { return nil }
func (httpAppStub) ValidatePkg(context.Context, int, domain.ValidatePkgRequest) ([]domain.Application, error) {
	return nil, nil
}
func (httpAppStub) GetAppConfigurations(context.Context, int, int) ([]domain.ApplicationConfigurationLink, error) {
	return nil, nil
}
func (httpAppStub) UpdateAppConfigurations(context.Context, int, domain.LinkConfigurationsToAppRequest) error {
	return nil
}
func (httpAppStub) GetVersionConfigurations(context.Context, int, int) ([]domain.ApplicationVersionConfigurationLink, error) {
	return nil, nil
}
func (httpAppStub) UpdateVersionConfigurations(context.Context, int, domain.LinkConfigurationsToAppVersionRequest) error {
	return nil
}
func (httpAppStub) AdminSearch(context.Context, string) ([]domain.Application, error) { return nil, nil }
func (httpAppStub) TurnIntoCommon(context.Context, int) error { return nil }

func intPtr(n int) *int { return &n }

func setupAppRouter(h *Handler, withPrincipal bool, perms ...string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if withPrincipal {
			platformauth.WithPrincipal(c, &platformauth.Principal{
				ID: 1, CustomerID: 1, AuthLoaded: true, Permissions: perms,
			})
		}
		c.Next()
	})
	h.Register(r.Group("/rest/private/applications"))
	return r
}

func TestSearch_forbiddenWithoutPrincipal(t *testing.T) {
	h := NewHandler(appapp.NewService(httpAppStub{}))
	r := setupAppRouter(h, false)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rest/private/applications/search", nil))
	if w.Code != http.StatusForbidden {
		t.Fatalf("status %d", w.Code)
	}
}

func TestSearch_ok(t *testing.T) {
	h := NewHandler(appapp.NewService(httpAppStub{}))
	r := setupAppRouter(h, true, platformauth.PermApplications)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rest/private/applications/search", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func TestListVersions_ok(t *testing.T) {
	h := NewHandler(appapp.NewService(httpAppStub{}))
	r := setupAppRouter(h, true, platformauth.PermApplications)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rest/private/applications/1/versions", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
}
