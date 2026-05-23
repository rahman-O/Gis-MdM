package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/devicelog/domain"
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
		SELECT id, customerid, logspreserveperiod FROM plugin_devicelog_settings WHERE customerid = $1`, customerID).
		Scan(&s.ID, &s.CustomerID, &s.LogsPreservePeriod)
	return s, err
}

func (r *Repository) SaveSettings(ctx context.Context, s domain.Settings) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE plugin_devicelog_settings SET logspreserveperiod = $1 WHERE customerid = $2`,
		s.LogsPreservePeriod, s.CustomerID)
	return err
}

func (r *Repository) UpsertRule(ctx context.Context, rule domain.Rule) (int64, error) {
	if rule.ID == 0 {
		return r.insertRule(ctx, rule)
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE plugin_devicelog_settings_rules
		SET name = $1, active = $2, applicationid = $3, severity = $4, filter = $5, groupid = $6, configurationid = $7
		WHERE id = $8`, rule.Name, rule.Active, rule.ApplicationID, rule.Severity, rule.Filter, rule.GroupID, rule.ConfigurationID, rule.ID)
	return rule.ID, err
}

func (r *Repository) insertRule(ctx context.Context, rule domain.Rule) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO plugin_devicelog_settings_rules (settingid, name, active, applicationid, severity, filter, groupid, configurationid)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		rule.SettingID, rule.Name, rule.Active, rule.ApplicationID, rule.Severity, rule.Filter, rule.GroupID, rule.ConfigurationID).Scan(&id)
	return id, err
}

func (r *Repository) DeleteRule(ctx context.Context, customerID, ruleID int64) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM plugin_devicelog_settings_rules r
		USING plugin_devicelog_settings s
		WHERE r.id = $1 AND r.settingid = s.id AND s.customerid = $2`, ruleID, customerID)
	return err
}

func (r *Repository) DeviceByNumber(ctx context.Context, number string) (id, customerID int64, err error) {
	err = r.db.QueryRowContext(ctx, `SELECT id, customerid FROM devices WHERE lower(number) = lower($1)`, number).Scan(&id, &customerID)
	return
}

func (r *Repository) DefaultApplicationID(ctx context.Context, customerID int64) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `SELECT id FROM applications WHERE customerid = $1 ORDER BY id LIMIT 1`, customerID).Scan(&id)
	return id, err
}

func (r *Repository) InsertLogs(ctx context.Context, customerID, deviceID, appID int64, rows []domain.UploadRecord) error {
	for _, row := range rows {
		ts := row.Timestamp
		if ts == 0 {
			ts = time.Now().UnixMilli()
		}
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO plugin_devicelog_log (createtime, customerid, deviceid, applicationid, severity, message, severityorder)
			VALUES ($1, $2, $3, $4, $5, $6, 0)`, ts, customerID, deviceID, appID, row.Severity, row.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) SearchLogs(ctx context.Context, customerID int64, f domain.LogFilter) ([]domain.LogRecord, int64, error) {
	if f.PageSize <= 0 {
		f.PageSize = 50
	}
	if f.PageNum <= 0 {
		f.PageNum = 1
	}
	offset := (f.PageNum - 1) * f.PageSize
	var total int64
	var rows *sql.Rows
	var err error
	if f.DeviceID > 0 {
		_ = r.db.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM plugin_devicelog_log WHERE customerid = $1 AND deviceid = $2`, customerID, f.DeviceID).Scan(&total)
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, createtime, deviceid, applicationid, COALESCE(severity,''), COALESCE(message,'')
			FROM plugin_devicelog_log WHERE customerid = $1 AND deviceid = $2
			ORDER BY createtime DESC LIMIT $3 OFFSET $4`, customerID, f.DeviceID, f.PageSize, offset)
	} else {
		_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM plugin_devicelog_log WHERE customerid = $1`, customerID).Scan(&total)
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, createtime, deviceid, applicationid, COALESCE(severity,''), COALESCE(message,'')
			FROM plugin_devicelog_log WHERE customerid = $1
			ORDER BY createtime DESC LIMIT $2 OFFSET $3`, customerID, f.PageSize, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []domain.LogRecord
	for rows.Next() {
		var l domain.LogRecord
		if err := rows.Scan(&l.ID, &l.CreateTime, &l.DeviceID, &l.ApplicationID, &l.Severity, &l.Message); err != nil {
			return nil, 0, err
		}
		items = append(items, l)
	}
	return items, total, rows.Err()
}
