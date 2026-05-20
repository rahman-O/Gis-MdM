package postgres

import (
	"context"
	"database/sql"

	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/port"
)

type ConfigRepository struct {
	db *sql.DB
}

func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

var _ port.ConfigByKey = (*ConfigRepository)(nil)

func (r *ConfigRepository) ConfigurationByQRKey(ctx context.Context, key string) (*port.QRConfig, error) {
	var cfg port.QRConfig
	var mainAppID sql.NullInt64
	err := r.db.QueryRowContext(ctx, `
		SELECT c.id, c.name, c.customerid, COALESCE(cu.filesdir, ''), c.mainappid
		FROM configurations c
		JOIN customers cu ON cu.id = c.customerid
		WHERE c.qrcodekey IS NOT NULL AND lower(c.qrcodekey) = lower($1)`, key).
		Scan(&cfg.ID, &cfg.Name, &cfg.CustomerID, &cfg.FilesDir, &mainAppID)
	if err != nil {
		return nil, err
	}
	if mainAppID.Valid {
		_ = r.db.QueryRowContext(ctx, `
			SELECT a.pkg, COALESCE(av.url, '')
			FROM applicationversions av
			JOIN applications a ON a.id = av.applicationid
			WHERE av.id = $1`, mainAppID.Int64).Scan(&cfg.MainAppPkg, &cfg.MainAppURL)
	}
	return &cfg, nil
}
