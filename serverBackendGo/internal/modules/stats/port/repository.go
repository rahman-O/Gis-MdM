package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/stats/domain"
)

// Repository persists usage statistics.
type Repository interface {
	Upsert(ctx context.Context, stats domain.UsageStats) error
}
