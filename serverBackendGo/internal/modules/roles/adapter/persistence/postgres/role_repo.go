package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/modules/roles/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/roles/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

var _ port.Repository = (*Repository)(nil)

func (r *Repository) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, COALESCE(description,''), superadmin
		FROM permissions WHERE superadmin = false ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Permission
	for rows.Next() {
		var p domain.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.SuperAdmin); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

const roleSelect = `
SELECT ur.id AS user_role_id, ur.name AS user_role_name, ur.description AS user_role_description,
       ur.superadmin AS user_role_super_admin,
       p.id AS permission_id, p.name AS permission_name, p.superadmin AS permission_super_admin
FROM userroles ur
LEFT JOIN userrolepermissions urp ON urp.roleid = ur.id
LEFT JOIN permissions p ON p.id = urp.permissionid
`

func (r *Repository) ListRoles(ctx context.Context, excludeOrgAdmin bool) ([]domain.Role, error) {
	q := roleSelect + ` WHERE NOT ur.superadmin`
	if excludeOrgAdmin {
		q += fmt.Sprintf(` AND ur.id <> %d`, platformauth.OrgAdminRoleID)
	}
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return aggregateRoles(rows)
}

func (r *Repository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	q := roleSelect + ` WHERE lower(ur.name)=lower($1)`
	rows, err := r.db.QueryContext(ctx, q, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	roles, err := aggregateRoles(rows)
	if err != nil || len(roles) == 0 {
		return nil, err
	}
	return &roles[0], nil
}

func aggregateRoles(rows *sql.Rows) ([]domain.Role, error) {
	byID := map[int]*domain.Role{}
	permSeen := map[int]map[int]struct{}{}
	for rows.Next() {
		var (
			roleID                            int
			roleName, roleDesc                sql.NullString
			roleSuper                         bool
			permID                            sql.NullInt64
			permName                          sql.NullString
			permSuper                         sql.NullBool
		)
		if err := rows.Scan(&roleID, &roleName, &roleDesc, &roleSuper, &permID, &permName, &permSuper); err != nil {
			return nil, err
		}
		role, ok := byID[roleID]
		if !ok {
			role = &domain.Role{
				ID: roleID, Name: roleName.String, Description: roleDesc.String, SuperAdmin: roleSuper,
			}
			byID[roleID] = role
			permSeen[roleID] = map[int]struct{}{}
		}
		if permID.Valid {
			pid := int(permID.Int64)
			if _, seen := permSeen[roleID][pid]; !seen {
				permSeen[roleID][pid] = struct{}{}
				role.Permissions = append(role.Permissions, domain.Permission{
					ID: pid, Name: permName.String, SuperAdmin: permSuper.Bool,
				})
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	out := make([]domain.Role, 0, len(byID))
	for _, role := range byID {
		out = append(out, *role)
	}
	return out, nil
}

func (r *Repository) Insert(ctx context.Context, role *domain.Role) error {
	var id int
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO userroles (name, description, superadmin) VALUES ($1,$2,false) RETURNING id`,
		role.Name, nullStr(role.Description),
	).Scan(&id)
	if err != nil {
		return err
	}
	role.ID = id
	return r.replacePermissions(ctx, id, role.Permissions)
}

func (r *Repository) Update(ctx context.Context, role *domain.Role) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE userroles SET name=$1, description=$2 WHERE id=$3`,
		role.Name, nullStr(role.Description), role.ID,
	)
	if err != nil {
		return err
	}
	return r.replacePermissions(ctx, role.ID, role.Permissions)
}

func (r *Repository) replacePermissions(ctx context.Context, roleID int, perms []domain.Permission) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM userrolepermissions WHERE roleid=$1`, roleID)
	if err != nil {
		return err
	}
	for _, p := range perms {
		if p.ID <= 0 {
			continue
		}
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO userrolepermissions (roleid, permissionid) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
			roleID, p.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, roleID int) error {
	if roleID == platformauth.OrgAdminRoleID {
		return fmt.Errorf("cannot delete org admin role")
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM userroles WHERE id=$1 AND superadmin=false`, roleID)
	return err
}

func nullStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}
