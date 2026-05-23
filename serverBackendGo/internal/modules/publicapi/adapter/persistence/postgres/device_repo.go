package postgres

import (
	"context"
	"database/sql"

	"github.com/gis-mdm/server-backend-go/internal/modules/publicapi/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/publicapi/port"
)

type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

var _ port.DeviceRepository = (*DeviceRepository)(nil)

func (r *DeviceRepository) FindDeviceByNumber(ctx context.Context, number string) (*domain.DeviceRef, error) {
	var d domain.DeviceRef
	err := r.db.QueryRowContext(ctx, `
		SELECT id, customerid FROM devices WHERE lower(number) = lower($1) LIMIT 1`, number).
		Scan(&d.ID, &d.CustomerID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DeviceRepository) CustomerFilesDir(ctx context.Context, customerID int) (string, error) {
	var dir sql.NullString
	err := r.db.QueryRowContext(ctx, `SELECT filesdir FROM customers WHERE id = $1`, customerID).Scan(&dir)
	if err != nil {
		return "", err
	}
	if dir.Valid {
		return dir.String, nil
	}
	return "", nil
}

func (r *DeviceRepository) HasDuplicateApp(ctx context.Context, customerID int, pkg, version string) (bool, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM applications a
		INNER JOIN applicationversions av ON av.applicationid = a.id
		WHERE a.customerid = $1 AND lower(a.pkg) = lower($2) AND av.version = $3`,
		customerID, pkg, version).Scan(&n)
	return n > 0, err
}

func (r *DeviceRepository) InsertApplication(ctx context.Context, customerID int, name, pkg, version, url string, flags domain.UploadAppRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var appID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO applications (name, pkg, customerid, type, showicon, system, url, runafterinstall, runatboot, usekiosk)
		VALUES ($1,$2,$3,'app',$4,$5,$6,$7,$8,$9) RETURNING id`,
		name, pkg, customerID, flags.ShowIcon, flags.System, nullStr(url),
		flags.RunAfterInstall, flags.RunAtBoot, flags.UseKiosk).Scan(&appID)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO applicationversions (applicationid, version, versioncode, url)
		VALUES ($1,$2,0,$3)`, appID, version, nullStr(url))
	if err != nil {
		return err
	}
	return tx.Commit()
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
