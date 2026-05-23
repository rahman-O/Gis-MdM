package middleware

import (
	"net"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// IPFilter restricts access by allowlist from the given environment variable.
// Empty allowlist allows all clients (legacy default).
func IPFilter(envKey string) gin.HandlerFunc {
	allowlist := parseAllowlist(os.Getenv(envKey))
	return func(c *gin.Context) {
		if len(allowlist) == 0 {
			c.Next()
			return
		}
		ip := clientIP(c)
		if !isAllowed(ip, allowlist) {
			c.AbortWithStatusJSON(403, gin.H{"status": "ERROR", "message": "Forbidden"})
			return
		}
		c.Next()
	}
}

func parseAllowlist(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func clientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return host
}

func isAllowed(ip string, allowlist []string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	for _, entry := range allowlist {
		if _, cidr, err := net.ParseCIDR(entry); err == nil {
			if cidr.Contains(parsed) {
				return true
			}
			continue
		}
		if entry == ip {
			return true
		}
	}
	return false
}
