package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/domain"
)

// RouteRepository persists enrollment routes.
type RouteRepository interface {
	ListViews(ctx context.Context, customerID int) ([]domain.EnrollmentRouteView, error)
	GetViewByID(ctx context.Context, customerID, id int) (*domain.EnrollmentRouteView, error)
	Create(ctx context.Context, customerID int, req domain.CreateRequest, qrcodeKey string, resolved domain.ResolvedBootstrap, containerAck bool) (int, error)
	Update(ctx context.Context, customerID, id int, req domain.UpdateRequest, resolved *domain.ResolvedBootstrap, containerAck *bool) error
	Delete(ctx context.Context, customerID, id int) error
	DeleteImpact(ctx context.Context, customerID, routeID int) (*domain.EnrollmentDeleteImpact, error)
	TreeNodeBelongsToCustomer(ctx context.Context, customerID, nodeID int) (bool, error)
	NodePlacementKind(ctx context.Context, customerID, nodeID int) (string, error)
	ListTreeNodeOptions(ctx context.Context, customerID, heavyThreshold int) ([]domain.TreeNodeOption, error)
	ListBootstrapApps(ctx context.Context, customerID int) ([]domain.BootstrapAppOption, error)
	RecordQRViewed(ctx context.Context, routeID int) error
}
