package application

import (
	"context"
	"testing"

	authdomain "github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

func TestListUsers_permissionDenied(t *testing.T) {
	svc := NewService(&stubRepo{})
	p := &platformauth.Principal{ID: 1, CustomerID: 1, RoleID: 3, AuthLoaded: true}
	_, err := svc.ListUsers(context.Background(), p, "")
	if err != ErrPermissionDenied {
		t.Fatalf("want permission denied, got %v", err)
	}
}

func TestListUsers_editableFlag(t *testing.T) {
	repo := &stubRepo{
		users: map[int64]*authdomain.User{
			1: {ID: 1, CustomerID: 1, UserRole: &authdomain.UserRole{ID: 1, SuperAdmin: true}},
			2: {ID: 2, CustomerID: 1, UserRole: &authdomain.UserRole{ID: 3}},
		},
	}
	repo.list = []*authdomain.User{repo.users[1], repo.users[2]}
	svc := NewService(repo)
	p := &platformauth.Principal{ID: 1, CustomerID: 1, SuperAdmin: true, AuthLoaded: true}
	views, err := svc.ListUsers(context.Background(), p, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(views) != 2 || views[0].Editable || !views[1].Editable {
		t.Fatalf("editable flags: %+v", views)
	}
}
