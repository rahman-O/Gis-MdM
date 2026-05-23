package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/summary/domain"
)

// SummaryRepository loads dashboard device statistics.
type SummaryRepository interface {
	HasDevicesTable(ctx context.Context) (bool, error)
	GetDeviceStats(ctx context.Context, customerID int, userID int64) (*domain.DeviceStats, error)
}
