package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	q := userDataSelect + ` WHERE users.id = $1`
	return r.queryUserByID(ctx, q, id)
}

func (r *UserRepository) queryUserByID(ctx context.Context, query string, id int64) (*domain.User, error) {
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query user: %w", err)
	}
	defer rows.Close()
	return aggregateUserRows(rows)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.findByEmail(ctx, email)
}

func (r *UserRepository) FindByPasswordResetToken(ctx context.Context, token string) (*domain.User, error) {
	q := userDataSelect + ` WHERE users.passwordresettoken = $1`
	return r.queryUsers(ctx, q, token)
}

func (r *UserRepository) SetPasswordResetToken(ctx context.Context, userID int64, token string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET passwordresettoken=$1 WHERE id=$2`, token, userID)
	return err
}

func (r *UserRepository) SetNewPassword(ctx context.Context, userID int64, passwordHash string, clearReset bool) error {
	if clearReset {
		_, err := r.db.ExecContext(ctx,
			`UPDATE users SET password=$1, passwordreset=false, passwordresettoken=NULL WHERE id=$2`,
			passwordHash, userID)
		return err
	}
	_, err := r.db.ExecContext(ctx, `UPDATE users SET password=$1 WHERE id=$2`, passwordHash, userID)
	return err
}

func (r *UserRepository) EmailUsedByCustomer(ctx context.Context, email string) (bool, error) {
	var n int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM users WHERE lower(email)=lower($1)`, email).Scan(&n)
	return n > 0, err
}

func (r *UserRepository) CustomerNameExists(ctx context.Context, name string) (bool, error) {
	var n int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM customers WHERE lower(name)=lower($1)`, name).Scan(&n)
	return n > 0, err
}

func (r *UserRepository) InsertPendingSignup(ctx context.Context, email, language, token string, signupTime int64) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO pendingsignup (email, signuptime, language, token)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (email) DO UPDATE SET signuptime=EXCLUDED.signuptime, language=EXCLUDED.language, token=EXCLUDED.token`,
		email, signupTime, language, token)
	return err
}

func (r *UserRepository) GetPendingSignupByToken(ctx context.Context, token string) (*domain.PendingSignup, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, signuptime, language, token FROM pendingsignup WHERE token=$1`, token)
	var p domain.PendingSignup
	if err := row.Scan(&p.ID, &p.Email, &p.SignupTime, &p.Language, &p.Token); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *UserRepository) GetPendingSignupByEmail(ctx context.Context, email string) (*domain.PendingSignup, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, signuptime, language, token FROM pendingsignup WHERE lower(email)=lower($1)`, email)
	var p domain.PendingSignup
	if err := row.Scan(&p.ID, &p.Email, &p.SignupTime, &p.Language, &p.Token); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *UserRepository) DeletePendingSignup(ctx context.Context, email string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM pendingsignup WHERE lower(email)=lower($1)`, email)
	return err
}

func (r *UserRepository) SignupCreateCustomer(ctx context.Context, p domain.SignupComplete) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	desc := p.Description
	if p.Company != "" {
		if desc != "" {
			desc += "\n"
		}
		desc += p.Company
	}

	var customerID int
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO customers (name, description, master, prefix)
		VALUES ($1, $2, false, 'hmdm-') RETURNING id`,
		p.Name, desc,
	).Scan(&customerID); err != nil {
		return 0, err
	}

	hash := crypto.HashFromMd5(p.PasswordMD5)
	displayName := strings.TrimSpace(p.FirstName + " " + p.LastName)
	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (login, email, name, password, customerid, userroleid, authtoken)
		VALUES ($1, $2, $3, $4, $5, 2, $6)`,
		p.Name, p.Email, displayName, hash, customerID, crypto.GenerateAuthToken(),
	)
	if err != nil {
		return 0, fmt.Errorf("insert user: %w", err)
	}
	_, _ = tx.ExecContext(ctx, `INSERT INTO settings (customerid, twofactor) VALUES ($1, false)`, customerID)
	return customerID, tx.Commit()
}

func (r *UserRepository) GetTwoFactorSecret(ctx context.Context, userID int64) (string, error) {
	var s sql.NullString
	err := r.db.QueryRowContext(ctx, `SELECT twofactorsecret FROM users WHERE id=$1`, userID).Scan(&s)
	return s.String, err
}

func (r *UserRepository) SetTwoFactorSecret(ctx context.Context, userID int64, secret string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET twofactorsecret=$1, twofactoraccepted=false WHERE id=$2`, secret, userID)
	return err
}

func (r *UserRepository) SetTwoFactorAccepted(ctx context.Context, userID int64, accepted bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET twofactoraccepted=$1 WHERE id=$2`, accepted, userID)
	return err
}

func (r *UserRepository) ClearTwoFactor(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET twofactorsecret=NULL, twofactoraccepted=false WHERE id=$1`, userID)
	return err
}
