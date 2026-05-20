package postgres

import (
	"context"
	"database/sql"

	authdomain "github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/customers/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

// UserLookup implements port.UserLookup.
type UserLookup struct {
	db *sql.DB
}

func NewUserLookup(db *sql.DB) *UserLookup {
	return &UserLookup{db: db}
}

var _ port.UserLookup = (*UserLookup)(nil)

func (u *UserLookup) FindOrgAdmin(ctx context.Context, customerID int) (*authdomain.User, error) {
	row := u.db.QueryRowContext(ctx, `
		SELECT id, login, email, name, password, customerid, userroleid, authtoken, passwordresettoken,
		       passwordreset, alldevicesavailable, allconfigavailable
		FROM users WHERE customerid = $1 AND userroleid = $2 LIMIT 1`,
		customerID, platformauth.OrgAdminRoleID)
	return scanOrgUser(row)
}

func (u *UserLookup) FindByLogin(ctx context.Context, login string) (*authdomain.User, error) {
	row := u.db.QueryRowContext(ctx, `
		SELECT id, login, email, name, password, customerid, userroleid, authtoken, passwordresettoken,
		       passwordreset, alldevicesavailable, allconfigavailable
		FROM users WHERE lower(login) = lower($1) LIMIT 1`, login)
	return scanOrgUser(row)
}

func (u *UserLookup) FindByEmail(ctx context.Context, email string) (*authdomain.User, error) {
	row := u.db.QueryRowContext(ctx, `
		SELECT id, login, email, name, password, customerid, userroleid, authtoken, passwordresettoken,
		       passwordreset, alldevicesavailable, allconfigavailable
		FROM users WHERE lower(email) = lower($1) LIMIT 1`, email)
	return scanOrgUser(row)
}

func scanOrgUser(row *sql.Row) (*authdomain.User, error) {
	var user authdomain.User
	var email, name, authToken, resetToken sql.NullString
	var roleID int
	if err := row.Scan(
		&user.ID, &user.Login, &email, &name, &user.Password, &user.CustomerID, &roleID,
		&authToken, &resetToken, &user.PasswordReset, &user.AllDevicesAvailable, &user.AllConfigAvailable,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if email.Valid {
		user.Email = email.String
	}
	if name.Valid {
		user.Name = name.String
	}
	if authToken.Valid {
		user.AuthToken = authToken.String
	}
	if resetToken.Valid {
		user.PasswordResetToken = resetToken.String
	}
	user.UserRole = &authdomain.UserRole{ID: roleID}
	return &user, nil
}

func (u *UserLookup) EnsureAuthToken(ctx context.Context, userID int64) (string, error) {
	var token, resetToken sql.NullString
	err := u.db.QueryRowContext(ctx,
		`SELECT authtoken, passwordresettoken FROM users WHERE id = $1`, userID).Scan(&token, &resetToken)
	if err != nil {
		return "", err
	}
	if resetToken.Valid && resetToken.String != "" {
		return "", port.ErrImpersonationBlocked
	}
	if token.Valid && token.String != "" {
		return token.String, nil
	}
	newTok := crypto.GenerateAuthToken()
	_, err = u.db.ExecContext(ctx, `UPDATE users SET authtoken = $1 WHERE id = $2`, newTok, userID)
	return newTok, err
}

func (u *UserLookup) UpdateOrgAdminMainDetails(ctx context.Context, userID int64, login, name, email string) error {
	_, err := u.db.ExecContext(ctx,
		`UPDATE users SET login=$1, name=$2, email=$3 WHERE id=$4`, login, name, email, userID)
	return err
}

func (u *UserLookup) InsertOrgAdmin(ctx context.Context, customerID int, login, name, email, passwordHash, authToken string, passwordReset bool) error {
	_, err := u.db.ExecContext(ctx, `
		INSERT INTO users (login, email, name, password, customerid, userroleid, authtoken, passwordreset)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		login, nullStr(email), name, passwordHash, customerID, platformauth.OrgAdminRoleID, authToken, passwordReset)
	return err
}
