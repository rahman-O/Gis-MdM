package application

import (
	"context"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

// JWTResult is the JWT login response body.
type JWTResult struct {
	IDToken string `json:"id_token"`
}

// JWTLogin issues a JWT for API clients (password: MD5 hex from UI, or raw for tools).
func (s *Service) JWTLogin(ctx context.Context, login, password string) (*JWTResult, string, error) {
	if login == "" || password == "" {
		return nil, "", BadRequest{}
	}
	passwordMD5, err := s.ResolvePassword(password)
	if err != nil {
		return nil, "", Unauthorized{}
	}
	user, err := s.repo.FindByLoginOrEmail(ctx, login)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		time.Sleep(bruteForceDelay)
		return nil, "", Unauthorized{}
	}
	if user.LastLoginFail > time.Now().UnixMilli()-1000 {
		time.Sleep(bruteForceDelay)
		return nil, "", Unauthorized{}
	}
	if !crypto.PasswordMatch(passwordMD5, user.Password) {
		_ = s.repo.SetLoginFailTime(ctx, user.ID, time.Now().UnixMilli())
		time.Sleep(bruteForceDelay)
		return nil, "", Unauthorized{}
	}

	go func() {
		_ = s.repo.RecordCustomerLastLogin(context.Background(), user.CustomerID, time.Now().UnixMilli())
	}()

	if err := s.repo.EnsureAuthToken(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := s.jwt.CreateToken(user.Login, user.AuthToken, false)
	if err != nil {
		return nil, "", err
	}
	return &JWTResult{IDToken: token}, "Bearer " + token, nil
}
