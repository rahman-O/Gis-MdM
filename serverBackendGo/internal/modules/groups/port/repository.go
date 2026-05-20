package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/groups/domain"
)

// GroupRepository persists device groups.
type GroupRepository interface {
	ListByCustomer(ctx context.Context, customerID int) ([]domain.Group, error)
	ListByValue(ctx context.Context, customerID int, value string) ([]domain.Group, error)
	GetByName(ctx context.Context, customerID int, name string) (*domain.Group, error)
	CountDevicesInGroup(ctx context.Context, groupID int) (int64, error)
	Insert(ctx context.Context, customerID int, name string) (int, error)
	Update(ctx context.Context, customerID int, g domain.Group) error
	Delete(ctx context.Context, customerID int, id int) error
	GrantCreatorAccess(ctx context.Context, userID int64, groupID int) error
	UserHasAllDevices(ctx context.Context, userID int64) (bool, error)
}
