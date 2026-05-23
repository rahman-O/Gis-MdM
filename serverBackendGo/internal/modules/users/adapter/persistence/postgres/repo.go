package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	authpostgres "github.com/gis-mdm/server-backend-go/internal/modules/auth/adapter/persistence/postgres"
	authdomain "github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/users/port"
)

// Repository implements users port using the auth user aggregate queries.
type Repository struct {
	inner *authpostgres.UserRepository
	db    *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{inner: authpostgres.NewUserRepository(db), db: db}
}

var _ port.Repository = (*Repository)(nil)

func (r *Repository) FindByID(ctx context.Context, id int64) (*authdomain.User, error) {
	return r.inner.FindByID(ctx, id)
}

func (r *Repository) FindByLogin(ctx context.Context, login string) (*authdomain.User, error) {
	return r.inner.FindByLoginOrEmail(ctx, login)
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*authdomain.User, error) {
	return r.inner.FindByEmail(ctx, email)
}

func (r *Repository) IsSingleCustomer(ctx context.Context) (bool, error) {
	return r.inner.IsSingleCustomer(ctx)
}

func (r *Repository) GetCustomerSettings(ctx context.Context, customerID int) (*authdomain.CustomerSettings, error) {
	return r.inner.GetCustomerSettings(ctx, customerID)
}

func (r *Repository) PasswordResetEnabled(ctx context.Context, customerID int) (bool, error) {
	_ = customerID
	return false, nil
}

func (r *Repository) ListByCustomer(ctx context.Context, customerID int, filter string) ([]*authdomain.User, error) {
	q := authpostgres.UserDataSelect + ` WHERE users.customerid = $1`
	args := []any{customerID}
	if filter != "" {
		q += ` AND (LOWER(users.name) LIKE $2 OR LOWER(users.login) LIKE $2 OR LOWER(users.email) LIKE $2)`
		args = append(args, "%"+strings.ToLower(filter)+"%")
	}
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()
	return authpostgres.AggregateUsersFromRows(rows)
}

func (r *Repository) UpdateMainDetails(ctx context.Context, u *authdomain.User) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET name=$1, email=$2, login=$3, userroleid=$4,
			alldevicesavailable=$5, allconfigavailable=$6
		WHERE id=$7 AND customerid=$8`,
		u.Name, nullIfEmpty(u.Email), u.Login, u.UserRole.ID,
		u.AllDevicesAvailable, u.AllConfigAvailable, u.ID, u.CustomerID,
	)
	if err != nil {
		return err
	}
	_, _ = r.db.ExecContext(ctx, `DELETE FROM userdevicegroupsaccess WHERE userid=$1`, u.ID)
	if !u.AllDevicesAvailable && len(u.Groups) > 0 {
		for _, g := range u.Groups {
			_, _ = r.db.ExecContext(ctx,
				`INSERT INTO userdevicegroupsaccess (userid, groupid) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
				u.ID, g.ID)
		}
	}
	_, _ = r.db.ExecContext(ctx, `DELETE FROM userconfigurationaccess WHERE userid=$1`, u.ID)
	if !u.AllConfigAvailable && len(u.Configurations) > 0 {
		for _, cfg := range u.Configurations {
			_, _ = r.db.ExecContext(ctx,
				`INSERT INTO userconfigurationaccess (userid, configurationid) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
				u.ID, cfg.ID)
		}
	}
	return nil
}

func (r *Repository) UpdatePassword(ctx context.Context, userID int64, passwordHash, authToken string, clear2FA bool, passwordReset bool, resetToken *string) error {
	if clear2FA {
		_, err := r.db.ExecContext(ctx, `
			UPDATE users SET password=$1, authtoken=$2, twofactorsecret=NULL, twofactoraccepted=false,
				passwordreset=$3, passwordresettoken=$4 WHERE id=$5`,
			passwordHash, authToken, passwordReset, resetToken, userID)
		return err
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET password=$1, authtoken=$2, passwordreset=$3, passwordresettoken=$4 WHERE id=$5`,
		passwordHash, authToken, passwordReset, resetToken, userID)
	return err
}

func (r *Repository) Insert(ctx context.Context, u *authdomain.User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (login, email, name, password, customerid, userroleid, authtoken,
			alldevicesavailable, allconfigavailable, passwordreset, passwordresettoken)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		u.Login, nullIfEmpty(u.Email), u.Name, u.Password, u.CustomerID, u.UserRole.ID, u.AuthToken,
		u.AllDevicesAvailable, u.AllConfigAvailable, u.PasswordReset, nullString(u.PasswordResetToken),
	)
	if err != nil {
		return err
	}
	var id int64
	if err := r.db.QueryRowContext(ctx, `SELECT id FROM users WHERE lower(login)=lower($1)`, u.Login).Scan(&id); err != nil {
		return err
	}
	u.ID = id
	return r.UpdateMainDetails(ctx, u)
}

func (r *Repository) Delete(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id=$1`, userID)
	return err
}

func nullIfEmpty(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}
