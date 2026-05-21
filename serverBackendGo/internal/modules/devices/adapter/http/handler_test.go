package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	devapp "github.com/gis-mdm/server-backend-go/internal/modules/devices/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/devices/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/devices/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type httpDeviceStub struct{}

func (httpDeviceStub) LoadUserScope(context.Context, int64) (*port.UserScope, error) {
	return &port.UserScope{CustomerID: 1, AllDevicesAvailable: true}, nil
}
func (httpDeviceStub) Search(context.Context, port.UserScope, domain.SearchRequest) ([]domain.DeviceView, error) {
	return []domain.DeviceView{}, nil
}
func (httpDeviceStub) Count(context.Context, port.UserScope, domain.SearchRequest) (int64, error) {
	return 0, nil
}
func (httpDeviceStub) ListConfigurations(context.Context, int) (map[int]domain.ConfigurationView, error) {
	return map[int]domain.ConfigurationView{}, nil
}
func (httpDeviceStub) GetByNumber(context.Context, port.UserScope, string) (*domain.DeviceView, error) {
	return nil, nil
}
func (httpDeviceStub) GetByID(context.Context, int, int) (*domain.DeviceView, error) { return nil, nil }
func (httpDeviceStub) ExistsNumber(context.Context, int, string, int) (bool, error)     { return false, nil }
func (httpDeviceStub) CountDevices(context.Context, int) (int64, error)                   { return 0, nil }
func (httpDeviceStub) DeviceLimit(context.Context, int) (int, error)                      { return 0, nil }
func (httpDeviceStub) Insert(context.Context, int, domain.SaveDevice) (int, error)      { return 1, nil }
func (httpDeviceStub) Update(context.Context, int, domain.SaveDevice) error             { return nil }
func (httpDeviceStub) UpdateConfigurationBulk(context.Context, int, []int, int) error   { return nil }
func (httpDeviceStub) Delete(context.Context, int, int) error                             { return nil }
func (httpDeviceStub) DeleteBulk(context.Context, int, []int) error                       { return nil }
func (httpDeviceStub) UpdateGroupBulk(context.Context, int, domain.GroupBulkRequest) error {
	return nil
}
func (httpDeviceStub) Autocomplete(context.Context, port.UserScope, string, int) ([]domain.LookupItem, error) {
	return nil, nil
}
func (httpDeviceStub) UpdateDescription(context.Context, int, int, string) error { return nil }
func (httpDeviceStub) ListAppSettings(context.Context, int) ([]domain.AppSetting, error) {
	return nil, nil
}
func (httpDeviceStub) SaveAppSettings(context.Context, int, []domain.AppSetting) error { return nil }

func setupDevicesRouter(h *Handler, withPrincipal bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if withPrincipal {
			platformauth.WithPrincipal(c, &platformauth.Principal{ID: 1, CustomerID: 1, AuthLoaded: true})
		}
		c.Next()
	})
	h.Register(r.Group("/rest/private/devices"))
	return r
}

func TestSearch_okWithPrincipal(t *testing.T) {
	h := NewHandler(devapp.NewService(httpDeviceStub{}, nil))
	r := setupDevicesRouter(h, true)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/rest/private/devices/search", bytes.NewBufferString(`{"pageNum":1,"pageSize":50}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func TestSearch_statusFilterAccepted(t *testing.T) {
	h := NewHandler(devapp.NewService(httpDeviceStub{}, nil))
	r := setupDevicesRouter(h, true)
	w := httptest.NewRecorder()
	body := `{"pageNum":1,"pageSize":50,"status":"green","sortBy":"LAST_UPDATE","sortDir":"desc"}`
	req := httptest.NewRequest(http.MethodPost, "/rest/private/devices/search", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}
