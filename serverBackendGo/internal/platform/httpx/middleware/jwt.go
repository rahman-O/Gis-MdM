package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	platformjwt "github.com/gis-mdm/server-backend-go/internal/platform/jwt"
)

// JWTAuth validates Bearer JWT and sets principal on context.
func JWTAuth(provider *platformjwt.Provider, lookup platformauth.UserLookup) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := resolveBearer(c.GetHeader("Authorization"))
		if token == "" {
			c.Next()
			return
		}
		claims, err := provider.ParseToken(token)
		if err != nil {
			c.AbortWithStatus(401)
			return
		}
		principal, err := lookup.LookupByLogin(c.Request.Context(), claims.Login)
		if err != nil || principal == nil {
			c.AbortWithStatus(403)
			return
		}
		if principal.AuthToken == "" || principal.AuthToken != claims.AuthToken {
			c.AbortWithStatus(403)
			return
		}
		platformauth.WithPrincipal(c, principal)
		c.Next()
	}
}

func resolveBearer(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
