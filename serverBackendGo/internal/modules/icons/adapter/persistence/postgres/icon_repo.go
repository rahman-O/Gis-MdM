package postgres

import (
	"context"
	"database/sql"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/icons/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/icons/port"
)

// IconRepository implements port.IconRepository.
type IconRepository struct {
	db *sql.DB
}

func NewIconRepository(db *sql.DB) *IconRepository {
	return &IconRepository{db: db}
}

var _ port.IconRepository = (*IconRepository)(nil)

const iconSelect = `
	SELECT icons.id, icons.customerid, icons.name, icons.fileid,
	       CASE
	         WHEN f.description IS NOT NULL AND f.description <> '' THEN f.description
	         WHEN f.external THEN COALESCE(f.externalurl, '')
	         ELSE COALESCE(f.filepath, '')
	       END AS filename
	FROM icons
	LEFT JOIN uploadedfiles f ON icons.fileid = f.id`

func (r *IconRepository) List(ctx context.Context, customerID int, filter string) ([]domain.Icon, error) {
	q := iconSelect + ` WHERE icons.customerid = $1`
	args := []any{customerID}
	if strings.TrimSpace(filter) != "" {
		q += ` AND icons.name ILIKE $2`
		args = append(args, "%"+filter+"%")
	}
	q += ` ORDER BY icons.name`
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Icon
	for rows.Next() {
		var ic domain.Icon
		if err := rows.Scan(&ic.ID, &ic.CustomerID, &ic.Name, &ic.FileID, &ic.FileName); err != nil {
			return nil, err
		}
		out = append(out, ic)
	}
	return out, rows.Err()
}

func (r *IconRepository) Save(ctx context.Context, icon domain.Icon) (*domain.Icon, error) {
	if icon.ID == nil || *icon.ID == 0 {
		var id int
		err := r.db.QueryRowContext(ctx, `
			INSERT INTO icons (customerid, name, fileid) VALUES ($1,$2,$3) RETURNING id`,
			icon.CustomerID, icon.Name, icon.FileID).Scan(&id)
		if err != nil {
			return nil, err
		}
		icon.ID = &id
	} else {
		_, err := r.db.ExecContext(ctx, `
			UPDATE icons SET name=$3, fileid=$4 WHERE id=$1 AND customerid=$2`,
			*icon.ID, icon.CustomerID, icon.Name, icon.FileID)
		if err != nil {
			return nil, err
		}
	}
	list, err := r.List(ctx, icon.CustomerID, "")
	if err != nil {
		return &icon, nil
	}
	for _, row := range list {
		if icon.ID != nil && row.ID != nil && *row.ID == *icon.ID {
			return &row, nil
		}
	}
	return &icon, nil
}

func (r *IconRepository) Delete(ctx context.Context, customerID, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM icons WHERE id=$1 AND customerid=$2`, id, customerID)
	return err
}
