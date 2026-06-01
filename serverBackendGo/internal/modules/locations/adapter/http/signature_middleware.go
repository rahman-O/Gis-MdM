package http

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/config"
)

const signatureMaxAgeSec = 300 // 5 minutes

// NewSignatureMiddleware returns a Gin middleware that validates HMAC-SHA256 signatures
// on incoming requests using the X-Device-Signature, X-Device-Id, and X-Request-Timestamp headers.
func NewSignatureMiddleware(cfg config.Config, log *slog.Logger) gin.HandlerFunc {
	maxAge := signatureMaxAgeSec
	secret := cfg.HashSecret

	return func(c *gin.Context) {
		signature := c.GetHeader("X-Device-Signature")
		deviceID := c.GetHeader("X-Device-Id")
		timestampStr := c.GetHeader("X-Request-Timestamp")

		if signature == "" || deviceID == "" || timestampStr == "" {
			log.Warn("signature middleware: missing required headers",
				"deviceId", deviceID,
				"hasSignature", signature != "",
				"hasTimestamp", timestampStr != "",
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing required signature headers",
			})
			return
		}

		// Validate timestamp freshness
		ts, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			log.Warn("signature middleware: invalid timestamp", "deviceId", deviceID, "timestamp", timestampStr)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid request timestamp",
			})
			return
		}

		now := time.Now().Unix()
		if math.Abs(float64(now-ts)) > float64(maxAge) {
			log.Warn("signature middleware: timestamp expired",
				"deviceId", deviceID,
				"requestTs", ts,
				"serverTs", now,
				"maxAge", maxAge,
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "request timestamp expired",
			})
			return
		}

		// Read body for signature computation
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Warn("signature middleware: failed to read body", "deviceId", deviceID, "err", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "failed to read request body",
			})
			return
		}
		// Restore body for downstream handlers
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		// Compute HMAC-SHA256 over deviceId + timestamp + body
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(deviceID))
		mac.Write([]byte(timestampStr))
		mac.Write(body)
		expected := hex.EncodeToString(mac.Sum(nil))

		// Timing-safe comparison
		if subtle.ConstantTimeCompare([]byte(expected), []byte(signature)) != 1 {
			log.Warn("signature middleware: signature mismatch",
				"deviceId", deviceID,
				"remoteAddr", c.ClientIP(),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid signature",
			})
			return
		}

		// Store device ID in context for downstream handlers
		c.Set("deviceId", deviceID)
		c.Next()
	}
}
