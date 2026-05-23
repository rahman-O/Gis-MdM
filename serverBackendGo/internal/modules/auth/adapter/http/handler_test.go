package http

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/port"
	"github.com/gis-mdm/server-backend-go/internal/platform/email"
	"github.com/gis-mdm/server-backend-go/internal/platform/jwt"
	"github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

func testService(repo port.UserRepository) *application.Service {
	return application.NewService(
		repo,
		jwt.NewProvider(jwt.Config{Secret: "s"}),
		email.NewService(true, slog.New(slog.NewTextHandler(os.Stderr, nil))),
		nil,
		false,
		true,
	)
}

func setupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := cookie.NewStore([]byte("test-secret"))
	r.Use(sessions.Sessions("hmdm_session", store))
	r.GET("/rest/public/auth/options", h.Options)
	r.POST("/rest/public/auth/login", h.Login)
	r.POST("/rest/public/auth/logout", h.Logout)
	r.POST("/rest/public/jwt/login", h.JWTLogin)
	return r
}

func TestOptions(t *testing.T) {
	svc := testService(&port.StubRepository{})
	r := setupRouter(NewHandler(svc))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rest/public/auth/options", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
}

func TestLogin_invalidPassword(t *testing.T) {
	repo := &port.StubRepository{User: &domain.User{ID: 1, Login: "admin", Password: "bad"}}
	r := setupRouter(NewHandler(testService(repo)))
	body, _ := json.Marshal(LoginRequest{Login: "admin", Password: crypto.MD5UpperHex("wrong")})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/rest/public/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "ERROR" {
		t.Fatalf("expected ERROR status, got %v", resp["status"])
	}
}

func TestJWTLogin_success_rawPassword(t *testing.T) {
	md5 := crypto.MD5UpperHex("admin")
	hash := crypto.HashFromMd5(md5)
	repo := &port.StubRepository{User: &domain.User{ID: 1, Login: "admin", Password: hash, CustomerID: 1}}
	r := setupRouter(NewHandler(testService(repo)))
	body, _ := json.Marshal(LoginRequest{Login: "admin", Password: "admin"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/rest/public/jwt/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
}

func TestJWTLogin_success_md5Hex(t *testing.T) {
	md5 := crypto.MD5UpperHex("admin")
	hash := crypto.HashFromMd5(md5)
	repo := &port.StubRepository{User: &domain.User{ID: 1, Login: "admin", Password: hash, CustomerID: 1}}
	r := setupRouter(NewHandler(testService(repo)))
	body, _ := json.Marshal(LoginRequest{Login: "admin", Password: md5})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/rest/public/jwt/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
	if w.Header().Get("Authorization") == "" {
		t.Fatal("expected Authorization header")
	}
}

func TestLogout_noContent(t *testing.T) {
	r := setupRouter(NewHandler(testService(&port.StubRepository{})))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/rest/public/auth/logout", nil))
	if w.Code != http.StatusNoContent {
		t.Fatalf("status %d", w.Code)
	}
}
