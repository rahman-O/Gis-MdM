package application

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/gis-mdm/server-backend-go/internal/modules/auth/port"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// TwoFactorService implements TOTP enrollment (legacy Headwind /rest/private/twofactor/*).
type TwoFactorService struct {
	repo   port.UserRepository
	issuer string
}

// NewTwoFactorService creates a two-factor service.
func NewTwoFactorService(repo port.UserRepository, issuer string) *TwoFactorService {
	if issuer == "" {
		issuer = "Headwind MDM"
	}
	return &TwoFactorService{repo: repo, issuer: issuer}
}

func (s *TwoFactorService) ensureKey(ctx context.Context, userID int64, accountName string) (*otp.Key, error) {
	secret, err := s.repo.GetTwoFactorSecret(ctx, userID)
	if err != nil {
		return nil, err
	}
	if secret != "" {
		return otp.NewKeyFromURL(s.otpURL(accountName, secret))
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuer,
		AccountName: accountName,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, err
	}
	if err := s.repo.SetTwoFactorSecret(ctx, userID, key.Secret()); err != nil {
		return nil, err
	}
	return key, nil
}

func (s *TwoFactorService) otpURL(accountName, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&period=30&digits=6",
		url.PathEscape(s.issuer), url.PathEscape(accountName), secret, url.QueryEscape(s.issuer))
}

// QRCodePNG returns a PNG QR for enrolling the authenticator app.
func (s *TwoFactorService) QRCodePNG(ctx context.Context, userID int64, accountName string) ([]byte, error) {
	key, err := s.ensureKey(ctx, userID, accountName)
	if err != nil {
		return nil, err
	}
	return qrcode.Encode(key.URL(), qrcode.Medium, 256)
}

// Verify checks a TOTP code for the user.
func (s *TwoFactorService) Verify(ctx context.Context, userID int64, code string) error {
	secret, err := s.repo.GetTwoFactorSecret(ctx, userID)
	if err != nil {
		return err
	}
	if secret == "" || !totp.Validate(code, secret) {
		return ErrInvalidTOTP
	}
	return nil
}

// SetAccepted marks 2FA as accepted for the user.
func (s *TwoFactorService) SetAccepted(ctx context.Context, userID int64) error {
	return s.repo.SetTwoFactorAccepted(ctx, userID, true)
}

// Reset clears two-factor for the user.
func (s *TwoFactorService) Reset(ctx context.Context, userID int64) error {
	return s.repo.ClearTwoFactor(ctx, userID)
}

// ErrInvalidTOTP is returned when the code does not match.
var ErrInvalidTOTP = errors.New("error.permission.denied")
