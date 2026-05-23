package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	hintapp "github.com/gis-mdm/server-backend-go/internal/modules/hints/application"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type httpStubRepo struct{}

func (httpStubRepo) GetHistory(context.Context, int64) ([]string, error) {
	return []string{}, nil
}
func (httpStubRepo) MarkShown(context.Context, int64, string) error { return nil }
func (httpStubRepo) Enable(context.Context, int64) error           { return nil }
func (httpStubRepo) Disable(context.Context, int64) error          { return nil }

func setupHintsRouter(h *Handler, withPrincipal bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if withPrincipal {
			platformauth.WithPrincipal(c, &platformauth.Principal{ID: 1, AuthLoaded: true})
		}
		c.Next()
	})
	g := r.Group("/rest/private/hints")
	h.Register(g)
	return r
}

func TestGetHistory_forbiddenWithoutPrincipal(t *testing.T) {
	h := NewHandler(hintapp.NewService(httpStubRepo{}))
	r := setupHintsRouter(h, false)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rest/private/hints/history", nil))
	if w.Code != http.StatusForbidden {
		t.Fatalf("status %d, want 403", w.Code)
	}
}

func TestGetHistory_okWithPrincipal(t *testing.T) {
	h := NewHandler(hintapp.NewService(httpStubRepo{}))
	r := setupHintsRouter(h, true)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rest/private/hints/history", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func TestMarkShown_jsonStringBody(t *testing.T) {
	h := NewHandler(hintapp.NewService(httpStubRepo{}))
	r := setupHintsRouter(h, true)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/rest/private/hints/history", strings.NewReader(`"hint.step.1"`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}
