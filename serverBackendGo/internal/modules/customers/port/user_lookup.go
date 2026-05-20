package port

import (
	"context"
	"errors"

	authdomain "github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
)

// ErrImpersonationBlocked when org admin has an active password-reset token.
var ErrImpersonationBlocked = errors.New("impersonation.blocked")

// UserLookup loads org-admin users for impersonation and customer create/update.
type UserLookup interface {
	FindOrgAdmin(ctx context.Context, customerID int) (*authdomain.User, error)
	FindByLogin(ctx context.Context, login string) (*authdomain.User, error)
	FindByEmail(ctx context.Context, email string) (*authdomain.User, error)
	EnsureAuthToken(ctx context.Context, userID int64) (string, error)
	UpdateOrgAdminMainDetails(ctx context.Context, userID int64, login, name, email string) error
	InsertOrgAdmin(ctx context.Context, customerID int, login, name, email, passwordHash, authToken string, passwordReset bool) error
}
