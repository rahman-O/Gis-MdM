package application

import (
	"context"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

const bruteForceDelay = time.Second

// Login authenticates with session flow (password: MD5 hex from UI, or raw for tools).
func (s *Service) Login(ctx context.Context, login, password string) (*domain.UserView, bool, error) {
	if login == "" || password == "" {
		return nil, false, AuthFailure{}
	}
	passwordMD5, err := s.ResolvePassword(password)
	if err != nil {
		return nil, false, AuthFailure{}
	}
	user, err := s.repo.FindByLoginOrEmail(ctx, login)
	if err != nil {
		return nil, false, err
	}
	if user == nil {
		time.Sleep(bruteForceDelay)
		return nil, false, AuthFailure{}
	}
	if user.LastLoginFail > time.Now().UnixMilli()-1000 {
		time.Sleep(bruteForceDelay)
		return nil, false, AuthFailure{}
	}
	if !crypto.PasswordMatch(passwordMD5, user.Password) {
		_ = s.repo.SetLoginFailTime(ctx, user.ID, time.Now().UnixMilli())
		time.Sleep(bruteForceDelay)
		return nil, false, AuthFailure{}
	}
	view, err := s.completeLogin(ctx, user)
	if err != nil {
		return nil, false, err
	}
	pending := view.TwoFactor != nil && *view.TwoFactor && (view.TwoFactorAccepted == nil || !*view.TwoFactorAccepted)
	return view, pending, nil
}

func (s *Service) completeLogin(ctx context.Context, user *domain.User) (*domain.UserView, error) {
	go func() {
		_ = s.repo.RecordCustomerLastLogin(context.Background(), user.CustomerID, time.Now().UnixMilli())
	}()

	if err := s.repo.EnsureAuthToken(ctx, user); err != nil {
		return nil, err
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
