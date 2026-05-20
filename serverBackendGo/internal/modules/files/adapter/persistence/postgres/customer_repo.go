package postgres

import (
	"context"
	"database/sql"

	"github.com/gis-mdm/server-backend-go/internal/modules/files/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/files/port"
)

// CustomerRepository loads customer file metadata.
type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

var _ port.CustomerRepository = (*CustomerRepository)(nil)

func (r *CustomerRepository) GetMeta(ctx context.Context, customerID int) (*domain.CustomerMeta, error) {
	var m domain.CustomerMeta
	var filesDir sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT id, COALESCE(filesdir, ''), COALESCE(sizelimit, 0), master
		FROM customers WHERE id = $1`, customerID).
		Scan(&m.ID, &filesDir, &m.SizeLimit, &m.Master)
	if err != nil {
		return nil, err
	}
	if filesDir.Valid {
		m.FilesDir = filesDir.String
	}
	return &m, nil
}

func (r *CustomerRepository) CountCustomers(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers`).Scan(&n)
	return n, err
}
