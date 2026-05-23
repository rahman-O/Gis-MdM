package application

import (
	"context"
	"errors"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/port"
	"github.com/gis-mdm/server-backend-go/internal/platform/email"
	"github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

// PasswordResetService handles public password recovery flows.
type PasswordResetService struct {
	repo  port.UserRepository
	email *email.Service
}

// NewPasswordResetService builds password reset use cases.
func NewPasswordResetService(repo port.UserRepository, email *email.Service) *PasswordResetService {
	return &PasswordResetService{repo: repo, email: email}
}

// Recover initiates password reset email (always OK to avoid user enumeration).
func (s *PasswordResetService) Recover(ctx context.Context, username string) error {
	if !s.email.IsConfigured() {
		return nil
	}
	user, err := s.repo.FindByLoginOrEmail(ctx, username)
	if err != nil {
		return err
	}
	if user == nil || strings.TrimSpace(user.Email) == "" {
		return nil
	}
	token := user.PasswordResetToken
	if token == "" {
		token = crypto.GenerateAuthToken()
		if err := s.repo.SetPasswordResetToken(ctx, user.ID, token); err != nil {
			return err
		}
	}
	_ = s.email.Send(user.Email, "Password recovery", s.email.RecoveryBody(token))
	return nil
}

// ResetSettingsView is returned by GET /passwordReset/settings/:token.
type ResetSettingsView struct {
	SingleCustomer bool `json:"singleCustomer"`
}

// ResetSettingsByToken returns minimal settings for reset UI.
func (s *PasswordResetService) ResetSettingsByToken(ctx context.Context, token string) (*ResetSettingsView, error) {
	user, err := s.repo.FindByPasswordResetToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	single, _ := s.repo.IsSingleCustomer(ctx)
	return &ResetSettingsView{SingleCustomer: single}, nil
}

// ResetPassword updates password from reset token.
func (s *PasswordResetService) ResetPassword(ctx context.Context, token, newPasswordMD5 string) (*domain.UserView, error) {
	user, err := s.repo.FindByPasswordResetToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	md5 := crypto.NormalizeLoginPassword(newPasswordMD5)
	hash := crypto.HashFromMd5(md5)
	if err := s.repo.SetNewPassword(ctx, user.ID, hash, true); err != nil {
		return nil, err
	}
	user.PasswordReset = false
	user.PasswordResetToken = ""
	user.Password = ""
	return domain.NewUserView(user), nil
}

// ErrUserNotFound is returned when reset token or user is invalid.
var ErrUserNotFound = errors.New("error.user.not.found")
