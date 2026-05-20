package postgres

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/deviceinfo/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetSettings(ctx context.Context, customerID int64) (domain.Settings, error) {
	var s domain.Settings
	err := r.db.QueryRowContext(ctx, `
		SELECT id, customerid, datapreserveperiod, senddata, intervalmins
		FROM plugin_deviceinfo_settings WHERE customerid = $1`, customerID).
		Scan(&s.ID, &s.CustomerID, &s.DataPreservePeriod, &s.SendData, &s.IntervalMins)
	return s, err
}

func (r *Repository) SaveSettings(ctx context.Context, s domain.Settings) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE plugin_deviceinfo_settings
		SET datapreserveperiod = $1, senddata = $2, intervalmins = $3
		WHERE customerid = $4`, s.DataPreservePeriod, s.SendData, s.IntervalMins, s.CustomerID)
	return err
}

func (r *Repository) DeviceByNumber(ctx context.Context, number string) (id, customerID int64, err error) {
	err = r.db.QueryRowContext(ctx, `
		SELECT id, customerid FROM devices WHERE lower(number) = lower($1)`, number).Scan(&id, &customerID)
	return
}

func (r *Repository) SaveDynamic(ctx context.Context, deviceID, customerID int64, items []domain.DynamicInfo) error {
	ts := time.Now().UnixMilli()
	var battery *int
	for _, it := range items {
		if strings.EqualFold(it.Attribute, "battery") || strings.EqualFold(it.Attribute, "batteryLevel") {
			if v, e := strconv.Atoi(it.Value); e == nil {
				battery = &v
			}
		}
	}
	var recordID int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO plugin_deviceinfo_deviceparams (deviceid, customerid, ts) VALUES ($1, $2, $3) RETURNING id`,
		deviceID, customerID, ts).Scan(&recordID)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO plugin_deviceinfo_deviceparams_device (recordid, batterylevel)
		VALUES ($1, $2) ON CONFLICT (recordid) DO UPDATE SET batterylevel = EXCLUDED.batterylevel`,
		recordID, battery)
	return err
}

func (r *Repository) ListRecords(ctx context.Context, deviceID int64, limit int) ([]domain.ParamsRecord, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT p.id, p.ts, d.batterylevel
		FROM plugin_deviceinfo_deviceparams p
		LEFT JOIN plugin_deviceinfo_deviceparams_device d ON d.recordid = p.id
		WHERE p.deviceid = $1 ORDER BY p.ts DESC LIMIT $2`, deviceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ParamsRecord
	for rows.Next() {
		var rec domain.ParamsRecord
		var bat sql.NullInt64
		if err := rows.Scan(&rec.ID, &rec.Ts, &bat); err != nil {
			return nil, err
		}
		if bat.Valid {
			v := int(bat.Int64)
			rec.BatteryLevel = &v
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}
