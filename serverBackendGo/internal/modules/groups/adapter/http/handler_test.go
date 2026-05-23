package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	groupapp "github.com/gis-mdm/server-backend-go/internal/modules/groups/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/groups/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type httpGroupStub struct{}

func (httpGroupStub) ListByCustomer(context.Context, int) ([]domain.Group, error) {
	return []domain.Group{{ID: 1, Name: "General"}}, nil
}
func (httpGroupStub) ListByValue(context.Context, int, string) ([]domain.Group, error) {
	return nil, nil
}
func (httpGroupStub) GetByName(context.Context, int, string) (*domain.Group, error) { return nil, nil }
func (httpGroupStub) CountDevicesInGroup(context.Context, int) (int64, error)      { return 0, nil }
func (httpGroupStub) Insert(context.Context, int, string) (int, error)             { return 1, nil }
func (httpGroupStub) Update(context.Context, int, domain.Group) error              { return nil }
func (httpGroupStub) Delete(context.Context, int, int) error                         { return nil }
func (httpGroupStub) GrantCreatorAccess(context.Context, int64, int) error         { return nil }
func (httpGroupStub) UserHasAllDevices(context.Context, int64) (bool, error)       { return true, nil }

func setupGroupsRouter(h *Handler, withPrincipal bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if withPrincipal {
			platformauth.WithPrincipal(c, &platformauth.Principal{ID: 1, CustomerID: 1, AuthLoaded: true})
		}
		c.Next()
	})
	h.Register(r.Group("/rest/private/groups"))
	return r
}

func TestSearch_forbiddenWithoutPrincipal(t *testing.T) {
	h := NewHandler(groupapp.NewService(httpGroupStub{}))
	r := setupGroupsRouter(h, false)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rest/private/groups/search", nil))
	if w.Code != http.StatusForbidden {
		t.Fatalf("status %d", w.Code)
	}
}

func TestSearch_okWithPrincipal(t *testing.T) {
	h := NewHandler(groupapp.NewService(httpGroupStub{}))
	r := setupGroupsRouter(h, true)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rest/private/groups/search", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}
