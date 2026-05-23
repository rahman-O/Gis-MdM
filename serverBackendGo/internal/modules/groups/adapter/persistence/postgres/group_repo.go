package postgres

import (
	"context"
	"database/sql"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/groups/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/groups/port"
)

// GroupRepository implements port.GroupRepository.
type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

var _ port.GroupRepository = (*GroupRepository)(nil)

func (r *GroupRepository) ListByCustomer(ctx context.Context, customerID int) ([]domain.Group, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name FROM groups WHERE customerid = $1 ORDER BY lower(name)`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanGroups(rows)
}

func (r *GroupRepository) ListByValue(ctx context.Context, customerID int, value string) ([]domain.Group, error) {
	pattern := "%" + strings.TrimSpace(value) + "%"
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name FROM groups
		WHERE customerid = $1 AND name ILIKE $2
		ORDER BY lower(name)`, customerID, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanGroups(rows)
}

func (r *GroupRepository) GetByName(ctx context.Context, customerID int, name string) (*domain.Group, error) {
	var g domain.Group
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name FROM groups
		WHERE customerid = $1 AND lower(name) = lower($2)`, customerID, name).Scan(&g.ID, &g.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GroupRepository) CountDevicesInGroup(ctx context.Context, groupID int) (int64, error) {
	var n int64
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM devicegroups WHERE groupid = $1`, groupID).Scan(&n)
	return n, err
}

func (r *GroupRepository) Insert(ctx context.Context, customerID int, name string) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO groups (name, customerid) VALUES ($1, $2) RETURNING id`, name, customerID).Scan(&id)
	return id, err
}

func (r *GroupRepository) Update(ctx context.Context, customerID int, g domain.Group) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE groups SET name = $1 WHERE id = $2 AND customerid = $3`, g.Name, g.ID, customerID)
	return err
}

func (r *GroupRepository) Delete(ctx context.Context, customerID int, id int) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM groups WHERE id = $1 AND customerid = $2`, id, customerID)
	return err
}

func (r *GroupRepository) GrantCreatorAccess(ctx context.Context, userID int64, groupID int) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO userdevicegroupsaccess (userid, groupid) VALUES ($1, $2)
		ON CONFLICT DO NOTHING`, userID, groupID)
	return err
}

func (r *GroupRepository) UserHasAllDevices(ctx context.Context, userID int64) (bool, error) {
	var all bool
	err := r.db.QueryRowContext(ctx, `
		SELECT alldevicesavailable FROM users WHERE id = $1`, userID).Scan(&all)
	return all, err
}

func scanGroups(rows *sql.Rows) ([]domain.Group, error) {
	var out []domain.Group
	for rows.Next() {
		var g domain.Group
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}
