package middleware

import (
	"github.com/gin-gonic/gin"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// EnrichPrincipal loads role/permissions for the current principal when missing.
func EnrichPrincipal(enricher platformauth.PrincipalEnricher) gin.HandlerFunc {
	return func(c *gin.Context) {
		if enricher == nil {
			c.Next()
			return
		}
		p, ok := platformauth.PrincipalFromContext(c)
		if ok && p != nil && !p.AuthLoaded {
			_ = enricher.EnrichPrincipal(c.Request.Context(), p)
			platformauth.WithPrincipal(c, p)
		}
		c.Next()
	}
}
