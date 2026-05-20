package auth

import (
	"context"

	"github.com/gin-gonic/gin"
)

type contextKey struct{}

// Principal is the authenticated user attached to a request.
type Principal struct {
	ID               int64
	Login            string
	AuthToken        string
	CustomerID       int
	PasswordReset    bool
	TwoFactorPending bool
}

const (
	sessionCredentialsKey  = "credentials"
	sessionTwoFactorNeeded = "twofactor"
)

// SessionCredentialsKey matches AuthFilter.sessionCredentials.
func SessionCredentialsKey() string  { return sessionCredentialsKey }
func SessionTwoFactorKey() string     { return sessionTwoFactorNeeded }

// WithPrincipal stores principal on gin and request context.
func WithPrincipal(c *gin.Context, p *Principal) {
	c.Set("principal", p)
	ctx := context.WithValue(c.Request.Context(), contextKey{}, p)
	c.Request = c.Request.WithContext(ctx)
}

// PrincipalFromContext returns the authenticated principal if present.
func PrincipalFromContext(c *gin.Context) (*Principal, bool) {
	if v, ok := c.Get("principal"); ok {
		if p, ok := v.(*Principal); ok {
			return p, true
		}
	}
	if v := c.Request.Context().Value(contextKey{}); v != nil {
		if p, ok := v.(*Principal); ok {
			return p, true
		}
	}
	return nil, false
}
