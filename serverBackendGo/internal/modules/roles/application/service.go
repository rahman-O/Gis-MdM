package application

import (
	"context"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/roles/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/roles/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// Service implements roles admin use cases.
type Service struct {
	repo           port.Repository
	singleCustomer func(context.Context) (bool, error)
}

func NewService(repo port.Repository, singleCustomer func(context.Context) (bool, error)) *Service {
	return &Service{repo: repo, singleCustomer: singleCustomer}
}

var (
	ErrPermissionDenied = errors.New("error.permission.denied")
	ErrDuplicateRole    = errors.New("error.duplicate.role")
)

func (s *Service) single(ctx context.Context) bool {
	if s.singleCustomer == nil {
		return true
	}
	ok, _ := s.singleCustomer(ctx)
	return ok
}

// ListPermissions returns non-superadmin permissions.
func (s *Service) ListPermissions(ctx context.Context, p *platformauth.Principal) ([]domain.Permission, error) {
	if p == nil || !p.CanManageRoles(s.single(ctx)) {
		return nil, ErrPermissionDenied
	}
	return s.repo.ListPermissions(ctx)
}

// ListRoles returns roles for admin UI (excludes org admin role per Java).
func (s *Service) ListRoles(ctx context.Context, p *platformauth.Principal) ([]domain.Role, error) {
	if p == nil || !p.CanManageRoles(s.single(ctx)) {
		return nil, ErrPermissionDenied
	}
	return s.repo.ListRoles(ctx, true)
}

// UpsertRole creates or updates a role.
func (s *Service) UpsertRole(ctx context.Context, p *platformauth.Principal, payload domain.RolePayload) error {
	if p == nil || !p.CanManageRoles(s.single(ctx)) {
		return ErrPermissionDenied
	}
	existing, err := s.repo.FindByName(ctx, payload.Name)
	if err != nil {
		return err
	}
	if existing != nil && (payload.ID == nil || existing.ID != *payload.ID) {
		return ErrDuplicateRole
	}
	role := &domain.Role{
		Name:        payload.Name,
		Description: payload.Description,
		Permissions: payload.Permissions,
	}
	if payload.ID == nil {
		return s.repo.Insert(ctx, role)
	}
	role.ID = *payload.ID
	return s.repo.Update(ctx, role)
}

// DeleteRole removes a role.
func (s *Service) DeleteRole(ctx context.Context, p *platformauth.Principal, roleID int) error {
	if p == nil || !p.CanManageRoles(s.single(ctx)) {
		return ErrPermissionDenied
	}
	return s.repo.Delete(ctx, roleID)
}
