package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/device_tree/domain"
)

// TreeRepository persists device tree nodes.
type TreeRepository interface {
	EnsureRoot(ctx context.Context, customerID int) (rootID int, err error)
	ListByCustomer(ctx context.Context, customerID int) ([]domain.TreeNode, error)
	GetByID(ctx context.Context, customerID, id int) (*domain.TreeNode, error)
	Create(ctx context.Context, customerID int, req domain.CreateNodeRequest) (*domain.TreeNode, error)
	Update(ctx context.Context, customerID, id int, req domain.UpdateNodeRequest) (*domain.TreeNode, error)
	DeleteWithRelocation(ctx context.Context, customerID, id, targetNodeID int) error
	CountDevicesInSubtree(ctx context.Context, customerID, nodeID int) (int, error)
}
