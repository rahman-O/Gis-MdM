package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// RequireAuth allows JWT principal or session (legacy AuthFilter).
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := platformauth.PrincipalFromContext(c); ok {
			c.Next()
			return
		}
		sess := sessions.Default(c)
		if v := sess.Get(platformauth.SessionCredentialsKey()); v != nil {
			if p, ok := v.(*platformauth.Principal); ok && p != nil {
				if p.PasswordReset {
					c.AbortWithStatus(http.StatusForbidden)
					return
				}
				if p.TwoFactorPending && !strings.Contains(c.Request.URL.Path, "/twofactor") {
					sess.Delete(platformauth.SessionCredentialsKey())
					sess.Delete(platformauth.SessionTwoFactorKey())
					_ = sess.Save()
					c.AbortWithStatus(http.StatusForbidden)
					return
				}
				platformauth.WithPrincipal(c, p)
				c.Next()
				return
			}
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}

// SessionStore saves principal after login.
func SessionStore(c *gin.Context, p *platformauth.Principal, twoFactorPending bool) {
	sess := sessions.Default(c)
	p.TwoFactorPending = twoFactorPending
	sess.Set(platformauth.SessionCredentialsKey(), p)
	if twoFactorPending {
		sess.Set(platformauth.SessionTwoFactorKey(), true)
	} else {
		sess.Delete(platformauth.SessionTwoFactorKey())
	}
	_ = sess.Save()
}

// SessionClear removes session on logout.
func SessionClear(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Clear()
	_ = sess.Save()
}

// ClearTwoFactorPending marks 2FA complete in session.
func ClearTwoFactorPending(c *gin.Context) {
	sess := sessions.Default(c)
	if v := sess.Get(platformauth.SessionCredentialsKey()); v != nil {
		if p, ok := v.(*platformauth.Principal); ok && p != nil {
			p.TwoFactorPending = false
			sess.Set(platformauth.SessionCredentialsKey(), p)
		}
	}
	sess.Delete(platformauth.SessionTwoFactorKey())
	_ = sess.Save()
}
