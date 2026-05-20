package postgres

import (
	"context"
	"database/sql"
	"fmt"

	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// EnrichPrincipal loads role id, superadmin flag, and permission names for a principal.
func (r *UserRepository) EnrichPrincipal(ctx context.Context, p *platformauth.Principal) error {
	if r.db == nil || p == nil {
		return nil
	}
	var roleID int
	var super bool
	err := r.db.QueryRowContext(ctx, `
		SELECT u.userroleid, ur.superadmin
		FROM users u
		INNER JOIN userroles ur ON ur.id = u.userroleid
		WHERE u.id = $1`, p.ID,
	).Scan(&roleID, &super)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("enrich principal role: %w", err)
	}
	p.RoleID = roleID
	p.SuperAdmin = super
	p.Permissions = nil

	rows, err := r.db.QueryContext(ctx, `
		SELECT DISTINCT p.name
		FROM permissions p
		INNER JOIN userrolepermissions urp ON urp.permissionid = p.id
		WHERE urp.roleid = $1`, roleID)
	if err != nil {
		return fmt.Errorf("enrich principal permissions: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		if name != "" {
			p.Permissions = append(p.Permissions, name)
		}
	}
	p.AuthLoaded = true
	return rows.Err()
}
