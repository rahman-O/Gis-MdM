package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	authdomain "github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	custapp "github.com/gis-mdm/server-backend-go/internal/modules/customers/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/customers/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type httpStubRepo struct{}

func (httpStubRepo) Search(context.Context, domain.SearchRequest) ([]domain.Customer, error) {
	return []domain.Customer{}, nil
}
func (httpStubRepo) Count(context.Context, domain.SearchRequest) (int64, error) { return 0, nil }
func (httpStubRepo) GetByID(context.Context, int) (*domain.Customer, error)     { return nil, nil }
func (httpStubRepo) GetByName(context.Context, string) (*domain.Customer, error) {
	return nil, nil
}
func (httpStubRepo) GetByEmail(context.Context, string) (*domain.Customer, error) {
	return nil, nil
}
func (httpStubRepo) Insert(context.Context, *domain.Customer) (int, error) { return 0, nil }
func (httpStubRepo) Update(context.Context, *domain.Customer) error       { return nil }
func (httpStubRepo) Delete(context.Context, int) error                   { return nil }
func (httpStubRepo) PrefixUsed(context.Context, string) (bool, error)    { return false, nil }

type httpStubUsers struct{}

func (httpStubUsers) FindOrgAdmin(context.Context, int) (*authdomain.User, error) {
	return nil, nil
}
func (httpStubUsers) FindByLogin(context.Context, string) (*authdomain.User, error) {
	return nil, nil
}
func (httpStubUsers) FindByEmail(context.Context, string) (*authdomain.User, error) {
	return nil, nil
}
func (httpStubUsers) EnsureAuthToken(context.Context, int64) (string, error) {
	return "t", nil
}
func (httpStubUsers) UpdateOrgAdminMainDetails(context.Context, int64, string, string, string) error {
	return nil
}
func (httpStubUsers) InsertOrgAdmin(context.Context, int, string, string, string, string, string, bool) error {
	return nil
}

func setupCustomersRouter(h *Handler, super bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if super {
			platformauth.WithPrincipal(c, &platformauth.Principal{ID: 1, SuperAdmin: true, AuthLoaded: true})
		} else {
			platformauth.WithPrincipal(c, &platformauth.Principal{ID: 2, SuperAdmin: false, AuthLoaded: true})
		}
		c.Next()
	})
	g := r.Group("/rest/private/customers")
	h.Register(g)
	return r
}

func TestSearch_forbiddenNonSuperAdmin(t *testing.T) {
	h := NewHandler(custapp.NewService(httpStubRepo{}, httpStubUsers{}))
	r := setupCustomersRouter(h, false)
	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"currentPage":1,"pageSize":10}`)
	req := httptest.NewRequest(http.MethodPost, "/rest/private/customers/search", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte(`error.permission.denied`)) {
		t.Fatalf("body %s", w.Body.String())
	}
}

func TestSearch_okSuperAdmin(t *testing.T) {
	h := NewHandler(custapp.NewService(httpStubRepo{}, httpStubUsers{}))
	r := setupCustomersRouter(h, true)
	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"currentPage":1,"pageSize":10}`)
	req := httptest.NewRequest(http.MethodPost, "/rest/private/customers/search", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte(`"status":"OK"`)) {
		t.Fatalf("body %s", w.Body.String())
	}
}
