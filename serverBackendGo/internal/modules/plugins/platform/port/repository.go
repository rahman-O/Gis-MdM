package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/platform/domain"
)

type Repository interface {
	FindActive(ctx context.Context) ([]domain.Plugin, error)
	FindAvailableByCustomer(ctx context.Context, customerID int64) ([]domain.Plugin, error)
	FindRegistered(ctx context.Context) ([]domain.Plugin, error)
	SaveDisabled(ctx context.Context, customerID int64, pluginIDs []int64) error
}
