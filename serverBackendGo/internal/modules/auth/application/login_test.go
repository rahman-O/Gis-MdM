package application

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/port"
	"github.com/gis-mdm/server-backend-go/internal/platform/email"
	"github.com/gis-mdm/server-backend-go/internal/platform/jwt"
	"github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

func testService(repo port.UserRepository) *Service {
	return NewService(repo, jwt.NewProvider(jwt.Config{Secret: "s"}), email.NewService(false, slog.New(slog.NewTextHandler(os.Stderr, nil))), nil, false, false)
}

func TestLogin_success(t *testing.T) {
	md5 := crypto.MD5UpperHex("admin")
	hash := crypto.HashFromMd5(md5)
	repo := &port.StubRepository{User: &domain.User{
		ID: 1, Login: "admin", Password: hash, CustomerID: 1,
		UserRole: &domain.UserRole{ID: 1, Name: "Admin", SuperAdmin: true},
	}}
	svc := testService(repo)
	view, _, err := svc.Login(context.Background(), "admin", md5)
	if err != nil {
		t.Fatal(err)
	}
	if view.AuthToken != "tok" {
		t.Fatalf("expected auth token, got %q", view.AuthToken)
	}
}

func TestLogin_rawPassword_normalized(t *testing.T) {
	md5 := crypto.MD5UpperHex("admin")
	hash := crypto.HashFromMd5(md5)
	repo := &port.StubRepository{User: &domain.User{
		ID: 1, Login: "admin", Password: hash, CustomerID: 1,
		UserRole: &domain.UserRole{ID: 1, Name: "Admin", SuperAdmin: true},
	}}
	svc := testService(repo)
	view, _, err := svc.Login(context.Background(), "admin", "admin")
	if err != nil {
		t.Fatal(err)
	}
	if view.Login != "admin" {
		t.Fatalf("unexpected login %q", view.Login)
	}
}

func TestLogin_failure(t *testing.T) {
	repo := &port.StubRepository{User: &domain.User{ID: 1, Login: "admin", Password: "x"}}
	svc := testService(repo)
	_, _, err := svc.Login(context.Background(), "admin", crypto.MD5UpperHex("wrong"))
	if _, ok := err.(AuthFailure); !ok {
		t.Fatalf("expected AuthFailure, got %v", err)
	}
}

func TestLogin_twoFactorPending(t *testing.T) {
	md5 := crypto.MD5UpperHex("admin")
	hash := crypto.HashFromMd5(md5)
	repo := &port.StubRepository{User: &domain.User{
		ID: 1, Login: "admin", Password: hash, CustomerID: 1,
		UserRole: &domain.UserRole{ID: 1, Name: "Admin"},
	}}
	svc := testService(repo)
	// default settings have TwoFactor false — pending should be false
	_, pending, err := svc.Login(context.Background(), "admin", md5)
	if err != nil || pending {
		t.Fatalf("login err=%v pending=%v", err, pending)
	}
}
