package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	stats := domain.EmptyDeviceStats()
	nowMs := time.Now().UnixMilli()

	green, err := r.countByOnlineWindow(ctx, customerID, userID, nowMs-2*3600*1000, 0)
	if err != nil {
		return nil, err
	}
	yellow, err := r.countByOnlineWindow(ctx, customerID, userID, nowMs-4*3600*1000, nowMs-2*3600*1000)
	if err != nil {
		return nil, err
	}
	red, err := r.countByOnlineWindow(ctx, customerID, userID, 0, nowMs-4*3600*1000)
	if err != nil {
		return nil, err
	}
	stats.StatusSummary = []domain.ChartItem{
		{StringAttr: "green", Number: int(green)},
		{StringAttr: "yellow", Number: int(yellow)},
		{StringAttr: "red", Number: int(red)},
	}

	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM devices WHERE customerid = $1`, customerID).Scan(&stats.DevicesTotal); err != nil {
		return nil, err
	}
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM devices WHERE customerid = $1 AND enrolltime > 0`, customerID).
		Scan(&stats.DevicesEnrolled); err != nil {
		return nil, err
	}
	monthAgo := nowMs - 30*86400*1000
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM devices WHERE customerid = $1 AND enrolltime >= $2`, customerID, monthAgo).
		Scan(&stats.DevicesEnrolledLastMonth); err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *SummaryRepository) countByOnlineWindow(ctx context.Context, customerID int, userID int64, minUpdate, maxUpdate int64) (int64, error) {
	query := `
		SELECT COUNT(DISTINCT d.id)
		FROM devices d
		INNER JOIN users u ON u.id = $2
		LEFT JOIN devicegroups dg ON d.id = dg.deviceid
		LEFT JOIN groups g ON dg.groupid = g.id
		LEFT JOIN userdevicegroupsaccess access ON g.id = access.groupid AND access.userid = u.id
		WHERE d.customerid = $1
		AND (u.alldevicesavailable = TRUE OR access.groupid IS NOT NULL)`
	args := []any{customerID, userID}
	if minUpdate > 0 {
		query += ` AND d.lastupdate >= $3`
		args = append(args, minUpdate)
	}
	if maxUpdate > 0 {
		n := len(args) + 1
		query += fmt.Sprintf(` AND d.lastupdate < $%d`, n)
		args = append(args, maxUpdate)
	}
	var n int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&n)
	return n, err
}
