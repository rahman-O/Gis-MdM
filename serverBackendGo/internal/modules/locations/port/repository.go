package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/locations/domain"
)

// Repository defines the data access interface for location records.
type Repository interface {
	InsertBatch(ctx context.Context, records []domain.LocationRecord) error
	GetLastRecord(ctx context.Context, deviceID int) (*domain.LocationRecord, error)
	QueryHistory(ctx context.Context, deviceID int, from, to int64, limit int) ([]domain.LocationRecord, int, error)
	QueryArchives(ctx context.Context, deviceID int, from, to int64) ([]domain.LocationArchive, error)
	InsertArchives(ctx context.Context, archives []domain.LocationArchive) error
	DeleteByIDs(ctx context.Context, ids []int64) error
	GetRecordsOlderThan(ctx context.Context, cutoffMs int64, batchSize int) ([]domain.LocationRecord, error)
}
