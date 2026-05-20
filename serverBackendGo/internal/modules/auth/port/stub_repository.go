package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
)

// StubRepository implements UserRepository with no-op defaults for tests.
type StubRepository struct {
	User *domain.User
}

func (s *StubRepository) FindByLoginOrEmail(_ context.Context, _ string) (*domain.User, error) {
	return s.User, nil
}
func (s *StubRepository) FindByID(context.Context, int64) (*domain.User, error)       { return s.User, nil }
func (s *StubRepository) FindByEmail(context.Context, string) (*domain.User, error)   { return nil, nil }
func (s *StubRepository) FindByPasswordResetToken(context.Context, string) (*domain.User, error) {
	return nil, nil
}
func (s *StubRepository) SetLoginFailTime(context.Context, int64, int64) error { return nil }
func (s *StubRepository) EnsureAuthToken(_ context.Context, u *domain.User) error {
	if u != nil {
		u.AuthToken = "tok"
	}
	return nil
}
func (s *StubRepository) SetPasswordResetToken(context.Context, int64, string) error { return nil }
func (s *StubRepository) SetNewPassword(context.Context, int64, string, bool) error  { return nil }
func (s *StubRepository) RecordCustomerLastLogin(context.Context, int, int64) error  { return nil }
func (s *StubRepository) IsSingleCustomer(context.Context) (bool, error)            { return true, nil }
func (s *StubRepository) GetCustomerSettings(context.Context, int) (*domain.CustomerSettings, error) {
	return &domain.CustomerSettings{}, nil
}
func (s *StubRepository) EmailUsedByCustomer(context.Context, string) (bool, error) { return false, nil }
func (s *StubRepository) CustomerNameExists(context.Context, string) (bool, error)  { return false, nil }
func (s *StubRepository) InsertPendingSignup(context.Context, string, string, string, int64) error {
	return nil
}
func (s *StubRepository) GetPendingSignupByToken(context.Context, string) (*domain.PendingSignup, error) {
	return nil, nil
}
func (s *StubRepository) GetPendingSignupByEmail(context.Context, string) (*domain.PendingSignup, error) {
	return nil, nil
}
func (s *StubRepository) DeletePendingSignup(context.Context, string) error { return nil }
func (s *StubRepository) SignupCreateCustomer(context.Context, domain.SignupComplete) (int, error) {
	return 0, nil
}
func (s *StubRepository) GetTwoFactorSecret(context.Context, int64) (string, error) { return "", nil }
func (s *StubRepository) SetTwoFactorSecret(context.Context, int64, string) error   { return nil }
func (s *StubRepository) SetTwoFactorAccepted(context.Context, int64, bool) error   { return nil }
func (s *StubRepository) ClearTwoFactor(context.Context, int64) error               { return nil }
