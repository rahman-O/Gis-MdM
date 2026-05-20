package postgres

import (
	"context"
	"database/sql"
)

// ConfigListItem is a minimal configuration row for dropdowns.
type ConfigListItem struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
}

// ConfigRepository lists configurations per tenant.
type ConfigRepository struct {
	db *sql.DB
}

func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

func (r *ConfigRepository) ListByCustomer(ctx context.Context, customerID int) ([]ConfigListItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name FROM configurations WHERE customerid = $1 ORDER BY lower(name)`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ConfigListItem
	for rows.Next() {
		var item ConfigListItem
		var name string
		if err := rows.Scan(&item.ID, &name); err != nil {
			return nil, err
		}
		item.Name = &name
		out = append(out, item)
	}
	return out, rows.Err()
}
