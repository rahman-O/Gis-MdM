package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/files/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/files/port"
)

// FileRepository implements port.FileRepository.
type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

var _ port.FileRepository = (*FileRepository)(nil)

func (r *FileRepository) List(ctx context.Context, customerID int, filter string) ([]domain.UploadedFile, error) {
	q := `SELECT id, customerid, filepath, description, uploadtime, devicepath, external, externalurl, replacevariables
		FROM uploadedfiles WHERE customerid = $1`
	args := []any{customerID}
	if strings.TrimSpace(filter) != "" {
		q += ` AND (filepath ILIKE $2 OR description ILIKE $2 OR externalurl ILIKE $2)`
		args = append(args, "%"+filter+"%")
	}
	q += ` ORDER BY filepath`
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFiles(rows)
}

func scanFiles(rows *sql.Rows) ([]domain.UploadedFile, error) {
	var out []domain.UploadedFile
	for rows.Next() {
		var f domain.UploadedFile
		var id int
		var path, desc, dev, extURL sql.NullString
		if err := rows.Scan(&id, &f.CustomerID, &path, &desc, &f.UploadTime, &dev, &f.External, &extURL, &f.ReplaceVariables); err != nil {
			return nil, err
		}
		f.ID = &id
		if path.Valid {
			f.FilePath = path.String
		}
		if desc.Valid {
			f.Description = desc.String
		}
		if dev.Valid {
			f.DevicePath = dev.String
		}
		if extURL.Valid {
			f.ExternalURL = extURL.String
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *FileRepository) GetByID(ctx context.Context, customerID, fileID int) (*domain.UploadedFile, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, customerid, filepath, description, uploadtime, devicepath, external, externalurl, replacevariables
		FROM uploadedfiles WHERE id = $1 AND customerid = $2`, fileID, customerID)
	var f domain.UploadedFile
	var rowID int
	var path, desc, dev, extURL sql.NullString
	if err := row.Scan(&rowID, &f.CustomerID, &path, &desc, &f.UploadTime, &dev, &f.External, &extURL, &f.ReplaceVariables); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	f.ID = &rowID
	if path.Valid {
		f.FilePath = path.String
	}
	if desc.Valid {
		f.Description = desc.String
	}
	if dev.Valid {
		f.DevicePath = dev.String
	}
	if extURL.Valid {
		f.ExternalURL = extURL.String
	}
	return &f, nil
}

func (r *FileRepository) Insert(ctx context.Context, f *domain.UploadedFile) error {
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO uploadedfiles (customerid, filepath, description, uploadtime, devicepath, external, externalurl, replacevariables)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id`,
		f.CustomerID, nullStr(f.FilePath), nullStr(f.Description), f.UploadTime, nullStr(f.DevicePath),
		f.External, nullStr(f.ExternalURL), f.ReplaceVariables).
		Scan(&id)
	if err != nil {
		return err
	}
	f.ID = &id
	return nil
}

func (r *FileRepository) Update(ctx context.Context, f *domain.UploadedFile) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE uploadedfiles SET filepath=$3, description=$4, uploadtime=$5, devicepath=$6,
			external=$7, externalurl=$8, replacevariables=$9
		WHERE id=$1 AND customerid=$2`,
		f.ID, f.CustomerID, nullStr(f.FilePath), nullStr(f.Description), f.UploadTime,
		nullStr(f.DevicePath), f.External, nullStr(f.ExternalURL), f.ReplaceVariables)
	return err
}

func (r *FileRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM uploadedfiles WHERE id = $1`, id)
	return err
}

func (r *FileRepository) IsUsedByConfiguration(ctx context.Context, fileID int) (bool, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM configurationfiles WHERE fileid = $1`, fileID).Scan(&n)
	return n > 0, err
}

func (r *FileRepository) IsUsedByIcon(ctx context.Context, fileID int) (bool, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM icons WHERE fileid = $1`, fileID).Scan(&n)
	return n > 0, err
}

func (r *FileRepository) UsingConfigurationNames(ctx context.Context, customerID, fileID int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.name FROM configurationfiles cf
		INNER JOIN configurations c ON c.id = cf.configurationid
		WHERE cf.fileid = $1 AND c.customerid = $2`, fileID, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStrings(rows)
}

func (r *FileRepository) UsingIconNames(ctx context.Context, customerID, fileID int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT ic.name FROM icons ic WHERE ic.fileid = $1 AND ic.customerid = $2`, fileID, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanStrings(rows)
}

func scanStrings(rows *sql.Rows) ([]string, error) {
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *FileRepository) GetFileConfigurations(ctx context.Context, customerID, userID, fileID int) ([]domain.FileConfigurationLink, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT cf.id, c.id, c.name, c.customerid, $2,
		       (cf.id IS NOT NULL) AS upload
		FROM configurations c
		INNER JOIN users u ON u.id = $3
		LEFT JOIN userconfigurationaccess access ON c.id = access.configurationid AND access.userid = u.id
		LEFT JOIN configurationfiles cf ON c.id = cf.configurationid AND cf.fileid = $2
		WHERE c.customerid = $1
		  AND (u.allconfigavailable = TRUE OR access.configurationid IS NOT NULL)
		ORDER BY lower(c.name)`, customerID, fileID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.FileConfigurationLink
	for rows.Next() {
		var l domain.FileConfigurationLink
		var linkID sql.NullInt64
		var upload bool
		if err := rows.Scan(&linkID, &l.ConfigurationID, &l.ConfigurationName, &l.CustomerID, &l.FileID, &upload); err != nil {
			return nil, err
		}
		if linkID.Valid {
			id := int(linkID.Int64)
			l.ID = &id
		}
		l.Upload = upload
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *FileRepository) DeleteConfigurationFile(ctx context.Context, linkID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM configurationfiles WHERE id = $1`, linkID)
	return err
}

func (r *FileRepository) InsertConfigurationFile(ctx context.Context, configurationID, fileID int, devicePath string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO configurationfiles (configurationid, fileid, path, devicepath)
		VALUES ($1,$2,$3,$4)`, configurationID, fileID, devicePath, devicePath)
	return err
}

func (r *FileRepository) CountByPath(ctx context.Context, customerID int, id *int, filePath string) (int64, error) {
	var n int64
	q := `SELECT COUNT(*) FROM uploadedfiles WHERE customerid = $1 AND filepath = $2`
	args := []any{customerID, filePath}
	if id != nil {
		q += ` AND id <> $3`
		args = append(args, *id)
	}
	err := r.db.QueryRowContext(ctx, q, args...).Scan(&n)
	return n, err
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// ErrNotFound for missing rows.
var ErrNotFound = fmt.Errorf("not found")
