package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// RequireSuperAdmin returns the principal when the caller is a super administrator.
func RequireSuperAdmin(c *gin.Context) (*Principal, bool) {
	p, ok := PrincipalFromContext(c)
	if !ok || p == nil || !p.SuperAdmin {
		response.PermissionDenied(c)
		return nil, false
	}
	return p, true
}
