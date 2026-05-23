package application

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/roles/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type stubRoleRepo struct {
	byName map[string]*domain.Role
}

func (s *stubRoleRepo) ListPermissions(context.Context) ([]domain.Permission, error) {
	return nil, nil
}
func (s *stubRoleRepo) ListRoles(context.Context, bool) ([]domain.Role, error) { return nil, nil }
func (s *stubRoleRepo) FindByName(_ context.Context, name string) (*domain.Role, error) {
	return s.byName[name], nil
}
func (s *stubRoleRepo) Insert(context.Context, *domain.Role) error { return nil }
func (s *stubRoleRepo) Update(context.Context, *domain.Role) error { return nil }
func (s *stubRoleRepo) Delete(context.Context, int) error          { return nil }

func TestUpsertRole_duplicateName(t *testing.T) {
	repo := &stubRoleRepo{byName: map[string]*domain.Role{"X": {ID: 5, Name: "X"}}}
	svc := NewService(repo, func(context.Context) (bool, error) { return true, nil })
	p := &platformauth.Principal{SuperAdmin: true, AuthLoaded: true}
	err := svc.UpsertRole(context.Background(), p, domain.RolePayload{Name: "X"})
	if err != ErrDuplicateRole {
		t.Fatalf("want duplicate role, got %v", err)
	}
}

func TestListPermissions_denied(t *testing.T) {
	svc := NewService(&stubRoleRepo{}, func(context.Context) (bool, error) { return false, nil })
	_, err := svc.ListPermissions(context.Background(), &platformauth.Principal{RoleID: 3})
	if err != ErrPermissionDenied {
		t.Fatalf("want permission denied, got %v", err)
	}
}
