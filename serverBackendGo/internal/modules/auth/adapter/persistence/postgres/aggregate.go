package postgres

import (
	"database/sql"
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
)

// UserDataSelect is the shared user aggregate query (exported for users module).
const UserDataSelect = userDataSelect

// AggregateUsersFromRows builds one user per id from a multi-row result set.
func AggregateUsersFromRows(rows *sql.Rows) ([]*domain.User, error) {
	byID := map[int64]*domain.User{}
	permSeen := map[int64]map[int]struct{}{}
	groupSeen := map[int64]map[int]struct{}{}
	configSeen := map[int64]map[int]struct{}{}

	for rows.Next() {
		u, err := scanUserRow(rows)
		if err != nil {
			return nil, err
		}
		existing, ok := byID[u.ID]
		if !ok {
			byID[u.ID] = u
			permSeen[u.ID] = map[int]struct{}{}
			groupSeen[u.ID] = map[int]struct{}{}
			configSeen[u.ID] = map[int]struct{}{}
			continue
		}
		mergeUserRelations(existing, u, permSeen[u.ID], groupSeen[u.ID], configSeen[u.ID])
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	out := make([]*domain.User, 0, len(byID))
	for _, u := range byID {
		out = append(out, u)
	}
	return out, nil
}

func mergeUserRelations(dst, src *domain.User, permSeen, groupSeen, configSeen map[int]struct{}) {
	if dst.UserRole == nil && src.UserRole != nil {
		dst.UserRole = src.UserRole
	}
	if src.UserRole != nil && len(src.UserRole.Permissions) > 0 && dst.UserRole != nil {
		for _, p := range src.UserRole.Permissions {
			if _, ok := permSeen[p.ID]; !ok {
				permSeen[p.ID] = struct{}{}
				dst.UserRole.Permissions = append(dst.UserRole.Permissions, p)
			}
		}
	}
	for _, g := range src.Groups {
		if _, ok := groupSeen[g.ID]; !ok {
			groupSeen[g.ID] = struct{}{}
			dst.Groups = append(dst.Groups, g)
		}
	}
	for _, c := range src.Configurations {
		if _, ok := configSeen[c.ID]; !ok {
			configSeen[c.ID] = struct{}{}
			dst.Configurations = append(dst.Configurations, c)
		}
	}
}

func scanUserRow(rows *sql.Rows) (*domain.User, error) {
	var (
		id, customerID                                                          int64
		name, login, email, authToken, resetToken, password, twoFactorSecret    sql.NullString
		allDev, allCfg, passReset, twoFactorAccepted, master                      bool
		lastFail                                                                int64
		roleID                                                                  sql.NullInt64
		roleName                                                                sql.NullString
		roleSuper                                                               sql.NullBool
		permID                                                                  sql.NullInt64
		permName                                                                sql.NullString
		permSuper                                                               sql.NullBool
		groupID                                                                 sql.NullInt64
		groupName                                                               sql.NullString
		configID                                                                sql.NullInt64
		configName                                                              sql.NullString
	)
	if err := rows.Scan(
		&id, &name, &login, &email, &customerID,
		&allDev, &allCfg, &passReset,
		&authToken, &resetToken, &lastFail, &password,
		&twoFactorSecret, &twoFactorAccepted, &master,
		&roleID, &roleName, &roleSuper,
		&permID, &permName, &permSuper,
		&groupID, &groupName,
		&configID, &configName,
	); err != nil {
		return nil, fmt.Errorf("scan user row: %w", err)
	}
	u := &domain.User{
		ID:                  id,
		Login:               login.String,
		Email:               email.String,
		Name:                name.String,
		Password:            password.String,
		CustomerID:          int(customerID),
		MasterCustomer:      master,
		AllDevicesAvailable: allDev,
		AllConfigAvailable:  allCfg,
		PasswordReset:       passReset,
		AuthToken:           authToken.String,
		PasswordResetToken:  resetToken.String,
		LastLoginFail:       lastFail,
		TwoFactorAccepted:   twoFactorAccepted,
	}
	if roleID.Valid {
		u.UserRole = &domain.UserRole{
			ID:         int(roleID.Int64),
			Name:       roleName.String,
			SuperAdmin: roleSuper.Bool,
		}
	}
	if permID.Valid {
		u.UserRole.Permissions = []domain.Permission{{
			ID: int(permID.Int64), Name: permName.String, SuperAdmin: permSuper.Bool,
		}}
	}
	if groupID.Valid {
		u.Groups = []domain.LookupItem{{ID: int(groupID.Int64), Name: groupName.String}}
	}
	if configID.Valid {
		u.Configurations = []domain.LookupItem{{ID: int(configID.Int64), Name: configName.String}}
	}
	return u, nil
}
