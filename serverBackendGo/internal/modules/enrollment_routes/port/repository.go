package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/domain"
)

// RouteRepository persists enrollment routes.
type RouteRepository interface {
	List(ctx context.Context, customerID int) ([]domain.Route, error)
	GetByID(ctx context.Context, customerID, id int) (*domain.RouteDetail, error)
	Create(ctx context.Context, customerID int, req domain.CreateRequest, qrcodeKey string) (int, error)
	Update(ctx context.Context, customerID, id int, req domain.UpdateRequest, qrcodeKey *string) error
	IsPublishedProfileVersion(ctx context.Context, customerID, profileVersionID int) (bool, error)
	ListPublishedProfileVersions(ctx context.Context, customerID int) ([]domain.PublishedProfileVersion, error)
	TreeNodeBelongsToCustomer(ctx context.Context, customerID, nodeID int) (bool, error)
}
