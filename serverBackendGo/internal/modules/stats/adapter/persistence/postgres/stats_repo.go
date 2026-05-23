package postgres

import (
	"context"
	"database/sql"

	"github.com/gis-mdm/server-backend-go/internal/modules/stats/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/stats/port"
)

type StatsRepository struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

var _ port.Repository = (*StatsRepository)(nil)

func (r *StatsRepository) Upsert(ctx context.Context, s domain.UsageStats) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO usagestats (
			ts, instanceid, webversion, community,
			devicestotal, devicesonline, cputotal, cpuused, ramtotal, ramused,
			scheme, arch, os
		) VALUES (
			CURRENT_DATE, $1, $2, $3,
			$4, $5, $6, $7, $8, $9,
			$10, $11, $12
		)
		ON CONFLICT (ts, instanceid) DO UPDATE SET
			webversion = EXCLUDED.webversion,
			community = EXCLUDED.community,
			devicestotal = EXCLUDED.devicestotal,
			devicesonline = EXCLUDED.devicesonline,
			cputotal = EXCLUDED.cputotal,
			cpuused = EXCLUDED.cpuused,
			ramtotal = EXCLUDED.ramtotal,
			ramused = EXCLUDED.ramused,
			scheme = EXCLUDED.scheme,
			arch = EXCLUDED.arch,
			os = EXCLUDED.os`,
		nullStr(s.InstanceID), nullStr(s.WebVersion), s.Community,
		s.DevicesTotal, s.DevicesOnline, s.CPUTotal, s.CPUUsed, s.RAMTotal, s.RAMUsed,
		nullStr(s.Scheme), nullStr(s.Arch), nullStr(s.OS))
	return err
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
