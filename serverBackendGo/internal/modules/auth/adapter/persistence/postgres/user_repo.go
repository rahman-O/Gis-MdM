package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/port"
	"github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

const userDataSelect = `
SELECT users.id, users.name, users.login, users.email, users.customerId,
       users.allDevicesAvailable, users.allConfigAvailable, users.passwordReset,
       users.authToken, users.passwordResetToken, users.lastLoginFail, users.password,
       users.twoFactorSecret, users.twoFactorAccepted,
       customers.master AS masterCustomer,
       userRoles.id AS user_role_id, userRoles.name AS user_role_name,
       userRoles.superadmin AS user_role_super_admin,
       permissions.id AS permission_id, permissions.name AS permission_name,
       permissions.superadmin AS permission_super_admin,
       groups.id AS groupId, groups.name AS groupName,
       configurations.id AS configurationId, configurations.name AS configurationName
FROM users
INNER JOIN customers ON customers.id = users.customerId
INNER JOIN userRoles ON userRoles.id = users.userRoleId
LEFT JOIN userRolePermissions urp ON urp.roleid = userRoles.id
LEFT JOIN permissions ON urp.permissionid = permissions.id
LEFT JOIN userDeviceGroupsAccess ON users.id = userDeviceGroupsAccess.userId
LEFT JOIN groups ON userDeviceGroupsAccess.groupId = groups.id
LEFT JOIN userConfigurationAccess ON users.id = userConfigurationAccess.userId
LEFT JOIN configurations ON userConfigurationAccess.configurationId = configurations.id
`

// UserRepository implements port.UserRepository and auth.UserLookup.
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

var (
	_ port.UserRepository = (*UserRepository)(nil)
	_ auth.UserLookup     = (*UserRepository)(nil)
)

func (r *UserRepository) FindByLoginOrEmail(ctx context.Context, login string) (*domain.User, error) {
	if r.db == nil {
		return nil, errors.New("database not configured")
	}
	u, err := r.findByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if u != nil {
		return u, nil
	}
	return r.findByEmail(ctx, login)
}

func (r *UserRepository) findByLogin(ctx context.Context, login string) (*domain.User, error) {
	q := userDataSelect + ` WHERE lower(users.login) = lower($1)`
	return r.queryUsers(ctx, q, login)
}

func (r *UserRepository) findByEmail(ctx context.Context, email string) (*domain.User, error) {
	q := userDataSelect + ` WHERE lower(users.email) = lower($1)`
	return r.queryUsers(ctx, q, email)
}

func (r *UserRepository) queryUsers(ctx context.Context, query string, arg string) (*domain.User, error) {
	rows, err := r.db.QueryContext(ctx, query, arg)
	if err != nil {
		return nil, fmt.Errorf("query user: %w", err)
	}
	defer rows.Close()
	return aggregateUserRows(rows)
}

func aggregateUserRows(rows *sql.Rows) (*domain.User, error) {
	users, err := AggregateUsersFromRows(rows)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

func (r *UserRepository) SetLoginFailTime(ctx context.Context, userID int64, ts int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET lastLoginFail=$1 WHERE id=$2`, ts, userID)
	return err
}

func (r *UserRepository) EnsureAuthToken(ctx context.Context, user *domain.User) error {
	if user.AuthToken != "" {
		return nil
	}
	token := crypto.GenerateAuthToken()
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET authToken=$1, password=$2 WHERE id=$3`,
		token, user.Password, user.ID,
	)
	if err != nil {
		return err
	}
	user.AuthToken = token
	return nil
}

func (r *UserRepository) RecordCustomerLastLogin(ctx context.Context, customerID int, ts int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE customers SET lastLoginTime=$1 WHERE id=$2`, ts, customerID,
	)
	return err
}

func (r *UserRepository) IsSingleCustomer(ctx context.Context) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM customers WHERE id > 1 LIMIT 1)`,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return !exists, nil
}

func (r *UserRepository) GetCustomerSettings(ctx context.Context, customerID int) (*domain.CustomerSettings, error) {
	var twoFactor bool
	var idleLogout sql.NullInt64
	err := r.db.QueryRowContext(ctx,
		`SELECT twoFactor, idleLogout FROM settings WHERE customerId=$1 LIMIT 1`,
		customerID,
	).Scan(&twoFactor, &idleLogout)
	if err == sql.ErrNoRows {
		return &domain.CustomerSettings{}, nil
	}
	if err != nil {
		return nil, err
	}
	s := &domain.CustomerSettings{TwoFactor: twoFactor}
	if idleLogout.Valid {
		v := int(idleLogout.Int64)
		s.IdleLogout = &v
	}
	return s, nil
}

// Lookup adapter for JWT middleware.
func (r *UserRepository) LookupByLogin(ctx context.Context, login string) (*auth.Principal, error) {
	u, err := r.FindByLoginOrEmail(ctx, login)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	perms := make([]string, 0)
	if u.UserRole != nil {
		for _, perm := range u.UserRole.Permissions {
			if perm.Name != "" {
				perms = append(perms, perm.Name)
			}
		}
	}
	roleID := 0
	super := false
	if u.UserRole != nil {
		roleID = u.UserRole.ID
		super = u.UserRole.SuperAdmin
	}
	return &auth.Principal{
		ID: u.ID, Login: u.Login, AuthToken: u.AuthToken, CustomerID: u.CustomerID,
		RoleID: roleID, SuperAdmin: super, Permissions: perms, AuthLoaded: true,
	}, nil
}
