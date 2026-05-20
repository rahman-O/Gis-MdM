package application

import (
	"context"
	"database/sql"
)

// RoleRow is a minimal user role for settings UI.
type RoleRow struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// RolesService lists assignable roles.
type RolesService struct {
	db *sql.DB
}

func NewRolesService(db *sql.DB) *RolesService {
	return &RolesService{db: db}
}

func (s *RolesService) ListRoles(ctx context.Context) ([]RoleRow, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name FROM userroles ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []RoleRow
	for rows.Next() {
		var r RoleRow
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
