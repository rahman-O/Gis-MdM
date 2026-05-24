package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/port"
)

type ConfigRepository struct {
	db *sql.DB
}

func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

var _ port.ConfigByKey = (*ConfigRepository)(nil)

func (r *ConfigRepository) CountCustomers(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers`).Scan(&n)
	return n, err
}

func (r *ConfigRepository) ConfigurationByQRKey(ctx context.Context, key string) (*port.QRConfig, error) {
	cfg, err := r.routeByQRKey(ctx, key)
	if err == nil {
		return cfg, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return r.configurationByQRKeyLegacy(ctx, key)
}

func (r *ConfigRepository) routeByQRKey(ctx context.Context, key string) (*port.QRConfig, error) {
	var cfg port.QRConfig
	var mainAppID sql.NullInt64
	var settingsJSON []byte
	var customerName sql.NullString
	var defaultDeviceIDMode string
	err := r.db.QueryRowContext(ctx, `
		SELECT er.id, er.name, COALESCE(er.qrcodekey, ''), er.customerid,
		       COALESCE(cu.filesdir, ''), COALESCE(cu.name, ''),
		       COALESCE(er.mainappid, pv.mainappid), pv.settingsjson,
		       COALESCE(NULLIF(TRIM(er.default_device_id_mode), ''), 'imei')
		FROM enrollment_routes er
		JOIN customers cu ON cu.id = er.customerid
		LEFT JOIN profile_versions pv ON pv.id = er.profile_version_id
		WHERE er.qrcodekey IS NOT NULL AND lower(er.qrcodekey) = lower($1)`, key).
		Scan(&cfg.ID, &cfg.Name, &cfg.QRCodeKey, &cfg.CustomerID, &cfg.FilesDir, &customerName, &mainAppID, &settingsJSON, &defaultDeviceIDMode)
	if err != nil {
		return nil, err
	}
	return r.finishQRConfig(ctx, &cfg, mainAppID, settingsJSON, customerName, defaultDeviceIDMode)
}

func (r *ConfigRepository) configurationByQRKeyLegacy(ctx context.Context, key string) (*port.QRConfig, error) {
	var cfg port.QRConfig
	var mainAppID sql.NullInt64
	var settingsJSON []byte
	var customerName sql.NullString
	var defaultDeviceIDMode string
	err := r.db.QueryRowContext(ctx, `
		SELECT c.id, c.name, COALESCE(c.qrcodekey, ''), c.customerid,
		       COALESCE(cu.filesdir, ''), COALESCE(cu.name, ''), c.mainappid, c.settingsjson,
		       COALESCE(NULLIF(TRIM(c.default_device_id_mode), ''), 'imei')
		FROM configurations c
		JOIN customers cu ON cu.id = c.customerid
		WHERE c.qrcodekey IS NOT NULL AND lower(c.qrcodekey) = lower($1)`, key).
		Scan(&cfg.ID, &cfg.Name, &cfg.QRCodeKey, &cfg.CustomerID, &cfg.FilesDir, &customerName, &mainAppID, &settingsJSON, &defaultDeviceIDMode)
	if err != nil {
		return nil, err
	}
	return r.finishQRConfig(ctx, &cfg, mainAppID, settingsJSON, customerName, defaultDeviceIDMode)
}

func (r *ConfigRepository) finishQRConfig(ctx context.Context, cfg *port.QRConfig, mainAppID sql.NullInt64, settingsJSON []byte, customerName sql.NullString, defaultDeviceIDMode string) (*port.QRConfig, error) {
	if customerName.Valid {
		cfg.CustomerName = customerName.String
	}
	cfg.DefaultDeviceIDMode = defaultDeviceIDMode
	parseQRSettings(settingsJSON, cfg)
	if mainAppID.Valid {
		cfg.MainAppVersionID = mainAppID.Int64
		var apkHash sql.NullString
		var appURL sql.NullString
		_ = r.db.QueryRowContext(ctx, `
			SELECT a.pkg, COALESCE(av.url, ''), COALESCE(av.filepath, ''), COALESCE(a.url, ''),
			       COALESCE(av.apkhash, '')
			FROM applicationversions av
			JOIN applications a ON a.id = av.applicationid
			WHERE av.id = $1`, mainAppID.Int64).Scan(
			&cfg.MainAppPkg, &cfg.MainAppURL, &cfg.MainAppFilePath, &appURL, &apkHash)
		if appURL.Valid {
			cfg.AppLevelURL = appURL.String
		}
		if apkHash.Valid {
			cfg.ApkHash = apkHash.String
		}
	}
	return cfg, nil
}

func parseQRSettings(raw []byte, cfg *port.QRConfig) {
	if len(raw) == 0 {
		return
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return
	}
	cfg.LauncherURL = jsonString(m, "launcherUrl")
	cfg.EventReceivingComponent = jsonString(m, "eventReceivingComponent")
	cfg.AdminExtras = jsonString(m, "adminExtras")
	cfg.QRParameters = jsonString(m, "qrParameters")
	cfg.WifiSSID = jsonString(m, "wifiSSID")
	cfg.WifiPassword = jsonString(m, "wifiPassword")
	cfg.WifiSecurityType = jsonString(m, "wifiSecurityType")
	cfg.MobileEnrollment = jsonBool(m, "mobileEnrollment")
	if v, ok := m["encryptDevice"]; ok {
		var b bool
		if json.Unmarshal(v, &b) == nil {
			cfg.EncryptDevice = b
		}
	}
}

func jsonString(m map[string]json.RawMessage, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return ""
	}
	return s
}

func jsonBool(m map[string]json.RawMessage, key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	var b bool
	_ = json.Unmarshal(v, &b)
	return b
}
