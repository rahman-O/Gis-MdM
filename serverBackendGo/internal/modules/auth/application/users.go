package application

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/port"
)

// UsersService exposes current-user reads for the React shell.
type UsersService struct {
	repo port.UserRepository
}

// NewUsersService builds users use cases.
func NewUsersService(repo port.UserRepository) *UsersService {
	return &UsersService{repo: repo}
}

// CurrentUser returns the authenticated user as UserView (password cleared).
func (s *UsersService) CurrentUser(ctx context.Context, userID int64) (*domain.UserView, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	settings, _ := s.repo.GetCustomerSettings(ctx, user.CustomerID)
	if settings != nil {
		if settings.TwoFactor {
			user.TwoFactor = true
		}
		user.IdleLogout = settings.IdleLogout
	}
	single, _ := s.repo.IsSingleCustomer(ctx)
	user.SingleCustomer = single
	user.Password = ""
	return domain.NewUserView(user), nil
}
