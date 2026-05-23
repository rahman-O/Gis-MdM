package postgres

import (
	"context"
	"database/sql"

	filesdomain "github.com/gis-mdm/server-backend-go/internal/modules/files/domain"
)

// UploadedFileRepository inserts uploadedfiles rows for icon uploads.
type UploadedFileRepository struct {
	db *sql.DB
}

func NewUploadedFileRepository(db *sql.DB) *UploadedFileRepository {
	return &UploadedFileRepository{db: db}
}

func (r *UploadedFileRepository) CustomerFilesDir(ctx context.Context, customerID int) (string, error) {
	var filesDir sql.NullString
	err := r.db.QueryRowContext(ctx, `SELECT COALESCE(filesdir, '') FROM customers WHERE id = $1`, customerID).Scan(&filesDir)
	if err != nil {
		return "", err
	}
	if filesDir.Valid {
		return filesDir.String, nil
	}
	return "", nil
}

func (r *UploadedFileRepository) InsertUploadedFile(ctx context.Context, f *filesdomain.UploadedFile) (*filesdomain.UploadedFile, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO uploadedfiles (customerid, filepath, uploadtime)
		VALUES ($1, $2, $3) RETURNING id`,
		f.CustomerID, f.FilePath, f.UploadTime).Scan(&id)
	if err != nil {
		return nil, err
	}
	f.ID = &id
	return f, nil
}
