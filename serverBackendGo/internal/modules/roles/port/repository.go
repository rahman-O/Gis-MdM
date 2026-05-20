package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/roles/domain"
)

// Repository loads and mutates roles and permissions.
type Repository interface {
	ListPermissions(ctx context.Context) ([]domain.Permission, error)
	ListRoles(ctx context.Context, excludeOrgAdmin bool) ([]domain.Role, error)
	FindByName(ctx context.Context, name string) (*domain.Role, error)
	Insert(ctx context.Context, role *domain.Role) error
	Update(ctx context.Context, role *domain.Role) error
	Delete(ctx context.Context, roleID int) error
}
