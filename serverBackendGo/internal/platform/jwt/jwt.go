package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const claimToken = "token"

// Provider issues and validates HS512 JWT tokens (legacy Headwind TokenProvider).
type Provider struct {
	secret                     []byte
	validity                   time.Duration
	validityRemember           time.Duration
}

// Config for JWT provider.
type Config struct {
	Secret           string
	ValiditySeconds  int64
	RememberSeconds  int64
}

// NewProvider creates a JWT provider.
func NewProvider(cfg Config) *Provider {
	secret := cfg.Secret
	if secret == "" {
		secret = "change-me"
	}
	validity := time.Duration(cfg.ValiditySeconds) * time.Second
	if validity <= 0 {
		validity = 24 * time.Hour
	}
	remember := time.Duration(cfg.RememberSeconds) * time.Second
	if remember <= 0 {
		remember = 30 * 24 * time.Hour
	}
	return &Provider{
		secret:           []byte(secret),
		validity:         validity,
		validityRemember: remember,
	}
}

// Claims parsed from a valid token.
type Claims struct {
	Login     string
	AuthToken string
}

// CreateToken builds a signed JWT for the user.
func (p *Provider) CreateToken(login, authToken string, rememberMe bool) (string, error) {
	ttl := p.validity
	if rememberMe {
		ttl = p.validityRemember
	}
	now := time.Now()
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":   login,
		claimToken: authToken,
		"exp":   now.Add(ttl).Unix(),
		"iat":   now.Unix(),
	})
	return t.SignedString(p.secret)
}

// ParseToken validates and returns claims.
func (p *Provider) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return p.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	sub, _ := claims["sub"].(string)
	tok, _ := claims[claimToken].(string)
	if sub == "" || tok == "" {
		return nil, errors.New("missing claims")
	}
	return &Claims{Login: sub, AuthToken: tok}, nil
}

// ValidateToken returns true if the token is valid.
func (p *Provider) ValidateToken(tokenString string) bool {
	_, err := p.ParseToken(tokenString)
	return err == nil
}
