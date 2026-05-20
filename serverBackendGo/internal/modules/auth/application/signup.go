package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/port"
	"github.com/gis-mdm/server-backend-go/internal/platform/email"
	"github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

// SignupService handles customer self-registration.
type SignupService struct {
	repo           port.UserRepository
	email          *email.Service
	customerSignup bool
}

// NewSignupService builds signup use cases.
func NewSignupService(repo port.UserRepository, email *email.Service, customerSignup bool) *SignupService {
	return &SignupService{repo: repo, email: email, customerSignup: customerSignup}
}

// Enabled reports whether self-signup is available.
func (s *SignupService) Enabled() bool {
	return s.customerSignup && s.email.IsConfigured()
}

// VerifyEmail starts signup for an email address.
func (s *SignupService) VerifyEmail(ctx context.Context, emailAddr, language string) error {
	if !s.customerSignup || !s.email.IsConfigured() {
		return ErrSignupDisabled
	}
	emailAddr = strings.ToLower(strings.TrimSpace(emailAddr))
	if emailAddr == "" {
		return ErrSignupDisabled
	}
	used, err := s.repo.EmailUsedByCustomer(ctx, emailAddr)
	if err != nil {
		return err
	}
	if used {
		return ErrDuplicateEmail
	}
	u, err := s.repo.FindByEmail(ctx, emailAddr)
	if err != nil {
		return err
	}
	if u != nil {
		return ErrDuplicateEmail
	}
	now := time.Now().UnixMilli()
	if pending, _ := s.repo.GetPendingSignupByEmail(ctx, emailAddr); pending != nil && pending.SignupTime+60000 > now {
		return ErrDuplicateEmail
	}
	token := crypto.GenerateAuthToken()
	if err := s.repo.InsertPendingSignup(ctx, emailAddr, language, token, now); err != nil {
		return err
	}
	_ = s.email.Send(emailAddr, "Verify your email", s.email.VerifySignupBody(token))
	return nil
}

// VerifyToken returns pending signup by token.
func (s *SignupService) VerifyToken(ctx context.Context, token string) (*domain.PendingSignup, error) {
	p, err := s.repo.GetPendingSignupByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrSignupTokenNotFound
	}
	return p, nil
}

// Complete finishes registration.
func (s *SignupService) Complete(ctx context.Context, req domain.SignupComplete) error {
	if !s.customerSignup {
		return ErrSignupDisabled
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return ErrSignupDisabled
	}
	exists, err := s.repo.CustomerNameExists(ctx, name)
	if err != nil {
		return err
	}
	if exists {
		return ErrDuplicateCustomer
	}
	pending, err := s.repo.GetPendingSignupByToken(ctx, req.Token)
	if err != nil {
		return err
	}
	if pending == nil {
		return ErrSignupTokenNotFound
	}
	req.Email = pending.Email
	req.Name = name
	md5 := crypto.NormalizeLoginPassword(req.PasswordMD5)
	req.PasswordMD5 = md5
	if _, err := s.repo.SignupCreateCustomer(ctx, req); err != nil {
		return err
	}
	_ = s.repo.DeletePendingSignup(ctx, pending.Email)
	_ = s.email.Send(pending.Email, "Welcome", "Your MDM account has been created.")
	return nil
}

var (
	ErrSignupDisabled      = errors.New("signup.disabled")
	ErrDuplicateEmail      = errors.New("signup.email.used")
	ErrDuplicateCustomer   = errors.New("error.duplicate.customer.name")
	ErrSignupTokenNotFound = errors.New("error.notfound.object")
)
