package port

import (
	"context"

	authdomain "github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
)

// Repository loads and mutates users for the users module.
type Repository interface {
	FindByID(ctx context.Context, id int64) (*authdomain.User, error)
	FindByLogin(ctx context.Context, login string) (*authdomain.User, error)
	FindByEmail(ctx context.Context, email string) (*authdomain.User, error)
	ListByCustomer(ctx context.Context, customerID int, filter string) ([]*authdomain.User, error)
	UpdateMainDetails(ctx context.Context, u *authdomain.User) error
	UpdatePassword(ctx context.Context, userID int64, passwordHash, authToken string, clear2FA bool, passwordReset bool, resetToken *string) error
	Insert(ctx context.Context, u *authdomain.User) error
	Delete(ctx context.Context, userID int64) error
	IsSingleCustomer(ctx context.Context) (bool, error)
	GetCustomerSettings(ctx context.Context, customerID int) (*authdomain.CustomerSettings, error)
	PasswordResetEnabled(ctx context.Context, customerID int) (bool, error)
}
