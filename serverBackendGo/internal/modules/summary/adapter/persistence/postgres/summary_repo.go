package postgres

import (
	"context"
	"database/sql"

	"github.com/gis-mdm/server-backend-go/internal/modules/summary/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/summary/port"
)

// SummaryRepository implements port.SummaryRepository.
type SummaryRepository struct {
	db *sql.DB
}

func NewSummaryRepository(db *sql.DB) *SummaryRepository {
	return &SummaryRepository{db: db}
}

var _ port.SummaryRepository = (*SummaryRepository)(nil)

func (r *SummaryRepository) HasDevicesTable(ctx context.Context) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = 'devices'
		)`).Scan(&exists)
	return exists, err
}

func (r *SummaryRepository) GetDeviceStats(ctx context.Context, customerID int, userID int64) (*domain.DeviceStats, error) {
	has, err := r.HasDevicesTable(ctx)
	if err != nil {
		return nil, err
	}
	if !has {
		return domain.EmptyDeviceStats(), nil
	}
	// Full SQL parity deferred to devices module; return empty shape so dashboard loads.
	_ = customerID
	_ = userID
	return domain.EmptyDeviceStats(), nil
}
