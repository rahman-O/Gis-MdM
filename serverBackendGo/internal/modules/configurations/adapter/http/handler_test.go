package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	cfgapp "github.com/gis-mdm/server-backend-go/internal/modules/configurations/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type httpCfgStub struct{}

func (httpCfgStub) ListByCustomer(context.Context, int) ([]domain.LookupItem, error) {
	name := "Default"
	return []domain.LookupItem{{ID: 1, Name: &name}}, nil
}
func (httpCfgStub) Search(context.Context, int) ([]domain.Configuration, error) {
	name := "Default"
	return []domain.Configuration{{ID: intPtr(1), Name: &name}}, nil
}
func (httpCfgStub) SearchByValue(context.Context, int, string) ([]domain.Configuration, error) {
	return nil, nil
}
func (httpCfgStub) GetByID(context.Context, int, int) (*domain.Configuration, error) {
	name := "Default"
	return &domain.Configuration{ID: intPtr(1), Name: &name}, nil
}
func (httpCfgStub) GetByName(context.Context, int, string) (*domain.Configuration, error) {
	return nil, nil
}
func (httpCfgStub) CountDevicesUsing(context.Context, int) (int64, error) { return 0, nil }
func (httpCfgStub) Insert(context.Context, int, domain.Configuration) (int, error) { return 1, nil }
func (httpCfgStub) Update(context.Context, int, domain.Configuration) error { return nil }
func (httpCfgStub) Delete(context.Context, int, int) error                     { return nil }
func (httpCfgStub) Copy(context.Context, int, domain.CopyRequest) (int, error) { return 2, nil }
func (httpCfgStub) ListAllApplicationsForPicker(context.Context, int) ([]domain.ConfigurationApplication, error) {
	return nil, nil
}
func (httpCfgStub) ListConfigurationApplications(context.Context, int, int) ([]domain.ConfigurationApplication, error) {
	return nil, nil
}
func (httpCfgStub) UpgradeApplication(context.Context, int, domain.UpgradeApplicationRequest) error {
	return nil
}

func intPtr(n int) *int { return &n }

func setupCfgRouter(h *Handler, withPrincipal bool, perms ...string) *gin.Engine {
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
	h.Register(r.Group("/rest/private/configurations"))
	return r
}

func TestSearch_forbiddenWithoutPrincipal(t *testing.T) {
	h := NewHandler(cfgapp.NewService(httpCfgStub{}, nil))
	r := setupCfgRouter(h, false)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rest/private/configurations/search", nil))
	if w.Code != http.StatusForbidden {
		t.Fatalf("status %d", w.Code)
	}
}

func TestSearch_okWithPermission(t *testing.T) {
	h := NewHandler(cfgapp.NewService(httpCfgStub{}, nil))
	r := setupCfgRouter(h, true, platformauth.PermConfigurations)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rest/private/configurations/search", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func TestGetByID_ok(t *testing.T) {
	h := NewHandler(cfgapp.NewService(httpCfgStub{}, nil))
	r := setupCfgRouter(h, true, platformauth.PermConfigurations)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rest/private/configurations/1", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
}

func TestList_okWithoutConfigurationsPermission(t *testing.T) {
	h := NewHandler(cfgapp.NewService(httpCfgStub{}, nil))
	r := setupCfgRouter(h, true)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rest/private/configurations/list", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
}
