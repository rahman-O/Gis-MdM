package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
)

// UserRepository loads and updates users for authentication and related flows.
type UserRepository interface {
	FindByLoginOrEmail(ctx context.Context, login string) (*domain.User, error)
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByPasswordResetToken(ctx context.Context, token string) (*domain.User, error)
	SetLoginFailTime(ctx context.Context, userID int64, ts int64) error
	EnsureAuthToken(ctx context.Context, user *domain.User) error
	SetPasswordResetToken(ctx context.Context, userID int64, token string) error
	SetNewPassword(ctx context.Context, userID int64, passwordHash string, clearReset bool) error
	RecordCustomerLastLogin(ctx context.Context, customerID int, ts int64) error
	IsSingleCustomer(ctx context.Context) (bool, error)
	GetCustomerSettings(ctx context.Context, customerID int) (*domain.CustomerSettings, error)
	EmailUsedByCustomer(ctx context.Context, email string) (bool, error)
	CustomerNameExists(ctx context.Context, name string) (bool, error)

	// Pending signup
	InsertPendingSignup(ctx context.Context, email, language, token string, signupTime int64) error
	GetPendingSignupByToken(ctx context.Context, token string) (*domain.PendingSignup, error)
	GetPendingSignupByEmail(ctx context.Context, email string) (*domain.PendingSignup, error)
	DeletePendingSignup(ctx context.Context, email string) error
	SignupCreateCustomer(ctx context.Context, p domain.SignupComplete) (int, error)

	// Two-factor
	GetTwoFactorSecret(ctx context.Context, userID int64) (string, error)
	SetTwoFactorSecret(ctx context.Context, userID int64, secret string) error
	SetTwoFactorAccepted(ctx context.Context, userID int64, accepted bool) error
	ClearTwoFactor(ctx context.Context, userID int64) error
}
