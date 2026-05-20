package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/sync/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/sync/port"
	"github.com/gis-mdm/server-backend-go/internal/platform/storage"
	sharedcrypto "github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

type DeviceSyncRepository struct {
	db *sql.DB
}

func NewDeviceSyncRepository(db *sql.DB) *DeviceSyncRepository {
	return &DeviceSyncRepository{db: db}
}

var _ port.SyncRepository = (*DeviceSyncRepository)(nil)

func scanDevice(row *sql.Row) (*domain.DeviceRecord, error) {
	var d domain.DeviceRecord
	var old sql.NullString
	var imei, phone, c1, c2, c3, info sql.NullString
	err := row.Scan(&d.ID, &d.CustomerID, &d.ConfigurationID, &d.Number, &old, &imei, &phone,
		&d.LastUpdate, &c1, &c2, &c3, &info)
	if err != nil {
		return nil, err
	}
	if old.Valid {
		d.OldNumber = &old.String
	}
	if imei.Valid {
		d.IMEI = &imei.String
	}
	if phone.Valid {
		d.Phone = &phone.String
	}
	if c1.Valid {
		d.Custom1 = &c1.String
	}
	if c2.Valid {
		d.Custom2 = &c2.String
	}
	if c3.Valid {
		d.Custom3 = &c3.String
	}
	if info.Valid {
		d.Info = &info.String
	}
	return &d, nil
}

const deviceSelect = `
	SELECT id, customerid, configurationid, number, oldnumber, imei, phone, lastupdate,
	       custom1, custom2, custom3, info
	FROM devices WHERE `

func (r *DeviceSyncRepository) FindByNumber(ctx context.Context, number string) (*domain.DeviceRecord, error) {
	return scanDevice(r.db.QueryRowContext(ctx, deviceSelect+` lower(number) = lower($1)`, number))
}

func (r *DeviceSyncRepository) FindByOldNumber(ctx context.Context, number string) (*domain.DeviceRecord, error) {
	return scanDevice(r.db.QueryRowContext(ctx, deviceSelect+` oldnumber IS NOT NULL AND lower(oldnumber) = lower($1)`, number))
}

func (r *DeviceSyncRepository) CountCustomers(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers`).Scan(&n)
	return n, err
}

func (r *DeviceSyncRepository) CreateOnDemand(ctx context.Context, number string, opts domain.DeviceCreateOptions, defaultCustomerID int64) (*domain.DeviceRecord, error) {
	configID, err := r.resolveConfigurationID(ctx, opts.Configuration, defaultCustomerID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UnixMilli()
	var id int64
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO devices (number, description, lastupdate, configurationid, customerid, enrolltime)
		VALUES ($1, '', 0, $2, $3, $4)
		RETURNING id`, number, configID, defaultCustomerID, now).Scan(&id)
	if err != nil {
		return nil, err
	}
	for _, gname := range opts.Groups {
		gname = strings.TrimSpace(gname)
		if gname == "" {
			continue
		}
		_, _ = r.db.ExecContext(ctx, `
			INSERT INTO devicegroups (deviceid, groupid)
			SELECT $1, g.id FROM groups g
			WHERE g.customerid = $2 AND lower(g.name) = lower($3)
			ON CONFLICT DO NOTHING`, id, defaultCustomerID, gname)
	}
	return r.FindByNumber(ctx, number)
}

func (r *DeviceSyncRepository) resolveConfigurationID(ctx context.Context, name string, customerID int64) (int64, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		var id int64
		err := r.db.QueryRowContext(ctx, `
			SELECT id FROM configurations WHERE customerid = $1 ORDER BY id LIMIT 1`, customerID).Scan(&id)
		return id, err
	}
	var id int64
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM configurations WHERE customerid = $1 AND lower(name) = lower($2)`, customerID, name).Scan(&id)
	return id, err
}

func (r *DeviceSyncRepository) CompleteMigration(ctx context.Context, deviceID int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE devices SET oldnumber = NULL WHERE id = $1`, deviceID)
	return err
}

func (r *DeviceSyncRepository) TouchLastUpdate(ctx context.Context, deviceID int64) error {
	now := time.Now().UnixMilli()
	_, err := r.db.ExecContext(ctx, `UPDATE devices SET lastupdate = $1 WHERE id = $2`, now, deviceID)
	return err
}

func (r *DeviceSyncRepository) UpdateInfo(ctx context.Context, deviceID int64, infoJSON, publicIP string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE devices SET info = $1, publicip = $2, lastupdate = $3 WHERE id = $4`,
		infoJSON, publicIP, time.Now().UnixMilli(), deviceID)
	return err
}

func (r *DeviceSyncRepository) UpdateCustomProps(ctx context.Context, deviceID int64, c1, c2, c3 *string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE devices SET custom1 = COALESCE($1, custom1), custom2 = COALESCE($2, custom2),
			custom3 = COALESCE($3, custom3) WHERE id = $4`, c1, c2, c3, deviceID)
	return err
}

func (r *DeviceSyncRepository) SaveApplicationSettings(ctx context.Context, deviceID int64, settings []domain.SyncApplicationSetting) error {
	for _, s := range settings {
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO deviceapplicationsettings (deviceid, applicationpkg, name, type, value)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT DO NOTHING`,
			deviceID, s.PackageID, s.Name, fmt.Sprintf("%d", s.Type), s.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *DeviceSyncRepository) BuildSyncResponse(ctx context.Context, dev domain.DeviceRecord, baseURL, filesDir, cpuArch, mobileName, vendor string) (*domain.SyncResponse, error) {
	var password, bg, fg sql.NullString
	var permissive bool
	var settingsJSON []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT password, backgroundcolor, textcolor, permissive, settingsjson
		FROM configurations WHERE id = $1`, dev.ConfigurationID).
		Scan(&password, &bg, &fg, &permissive, &settingsJSON)
	if err != nil {
		return nil, err
	}
	resp := &domain.SyncResponse{
		DeviceID:        dev.Number,
		ConfigurationID: dev.ConfigurationID,
		Applications:    []domain.SyncApplication{},
		Files:           []domain.SyncConfigurationFile{},
	}
	if password.Valid && password.String != "" {
		resp.Password = sharedcrypto.MD5UpperHex(password.String)
	}
	if bg.Valid {
		resp.BackgroundColor = &bg.String
	}
	if fg.Valid {
		resp.TextColor = &fg.String
	}
	resp.Permissive = &permissive
	if mobileName != "" {
		resp.AppName = &mobileName
	}
	if vendor != "" {
		resp.Vendor = &vendor
	}
	if dev.OldNumber != nil {
		resp.NewNumber = &dev.Number
	}
	resp.Custom1 = dev.Custom1
	resp.Custom2 = dev.Custom2
	resp.Custom3 = dev.Custom3

	var extra map[string]any
	_ = json.Unmarshal(settingsJSON, &extra)
	if v, ok := extra["pushOptions"].(string); ok {
		resp.PushOptions = &v
	}
	if v, ok := extra["requestUpdates"].(string); ok {
		resp.RequestUpdates = &v
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT a.id, a.name, a.pkg, COALESCE(av.version, ''), COALESCE(av.url, ''),
			COALESCE(a.type, 'app'), COALESCE(av.urlarm64, ''), COALESCE(av.urlarmeabi, ''), COALESCE(av.split, false)
		FROM configurationapplications ca
		JOIN applications a ON a.id = ca.applicationid
		LEFT JOIN applicationversions av ON av.id = ca.applicationversionid
		WHERE ca.configurationid = $1
		ORDER BY ca.screenorder NULLS LAST, a.name`, dev.ConfigurationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var app domain.SyncApplication
		var urlArm64, urlArm sql.NullString
		var split bool
		if err := rows.Scan(&app.ID, &app.Name, &app.Pkg, &app.Version, &app.URL, &app.Type, &urlArm64, &urlArm, &split); err != nil {
			return nil, err
		}
		if split {
			if strings.HasPrefix(cpuArch, "arm64") && urlArm64.Valid && urlArm64.String != "" {
				app.URL = urlArm64.String
			} else if urlArm.Valid && urlArm.String != "" {
				app.URL = urlArm.String
			}
		}
		resp.Applications = append(resp.Applications, app)
	}

	frows, err := r.db.QueryContext(ctx, `
		SELECT COALESCE(path, ''), COALESCE(externalurl, ''), COALESCE(url, ''), remove
		FROM configurationfiles WHERE configurationid = $1`, dev.ConfigurationID)
	if err != nil {
		return nil, err
	}
	defer frows.Close()
	for frows.Next() {
		var path, extURL, url string
		var remove bool
		if err := frows.Scan(&path, &extURL, &url, &remove); err != nil {
			return nil, err
		}
		f := domain.SyncConfigurationFile{Path: path, Remove: remove}
		if extURL != "" {
			f.URL = extURL
			f.External = true
		} else if path != "" {
			f.URL = storage.BuildPublicURL(baseURL, filesDir, path)
		} else if url != "" {
			f.URL = url
		}
		if f.Path == "" || (!remove && f.URL == "") {
			continue
		}
		resp.Files = append(resp.Files, f)
	}
	return resp, nil
}
