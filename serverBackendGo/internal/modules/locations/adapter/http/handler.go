package http

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/modules/locations/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/locations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/locations/port"
)

// LocationHandler handles REST endpoints for the locations module.
type LocationHandler struct {
	service       *application.LocationService
	repo          port.Repository
	log           *slog.Logger
	retentionDays int
}

// NewLocationHandler creates a new LocationHandler.
func NewLocationHandler(service *application.LocationService, repo port.Repository, log *slog.Logger, retentionDays int) *LocationHandler {
	if retentionDays <= 0 {
		retentionDays = 30
	}
	return &LocationHandler{
		service:       service,
		repo:          repo,
		log:           log,
		retentionDays: retentionDays,
	}
}

// BatchUpload handles POST /api/devices/:deviceId/locations/batch
// Parses the request body and delegates to the service layer.
func (h *LocationHandler) BatchUpload(c *gin.Context) {
	deviceIDStr := c.Param("deviceId")
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deviceId"})
		return
	}

	var entries []domain.BatchUploadEntry
	if err := c.ShouldBindJSON(&entries); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp := h.service.ProcessBatch(c.Request.Context(), deviceID, entries)
	c.JSON(http.StatusOK, resp)
}

// GetHistory handles GET /api/devices/:deviceId/locations
// Returns raw location records or routes to archive based on the time range.
func (h *LocationHandler) GetHistory(c *gin.Context) {
	deviceID, err := strconv.Atoi(c.Param("deviceId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deviceId"})
		return
	}

	from := parseQueryInt64(c, "from", 0)
	to := parseQueryInt64(c, "to", time.Now().UnixMilli())
	limit := parseQueryInt(c, "limit", 500)
	if limit > 5000 {
		limit = 5000
	}

	// If the requested range is entirely within the retention period, query raw records.
	// Otherwise, route to archives for the older portion.
	retentionCutoff := time.Now().Add(-time.Duration(h.retentionDays) * 24 * time.Hour).UnixMilli()

	if from >= retentionCutoff {
		// All data is in raw records
		records, total, err := h.repo.QueryHistory(c.Request.Context(), deviceID, from, to, limit)
		if err != nil {
			h.log.Error("query history failed", "deviceId", deviceID, "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"records": records,
			"total":   total,
			"source":  "raw",
		})
		return
	}

	// Mix: archives for old data, raw for recent
	archives, err := h.repo.QueryArchives(c.Request.Context(), deviceID, from, retentionCutoff)
	if err != nil {
		h.log.Error("query archives failed", "deviceId", deviceID, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	records, total, err := h.repo.QueryHistory(c.Request.Context(), deviceID, retentionCutoff, to, limit)
	if err != nil {
		h.log.Error("query history failed", "deviceId", deviceID, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"records":  records,
		"archives": archives,
		"total":    total,
		"source":   "mixed",
	})
}

// GetArchive handles GET /api/devices/:deviceId/locations/archive
// Returns archived hourly summaries for the given time range.
func (h *LocationHandler) GetArchive(c *gin.Context) {
	deviceID, err := strconv.Atoi(c.Param("deviceId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deviceId"})
		return
	}

	from := parseQueryInt64(c, "from", 0)
	to := parseQueryInt64(c, "to", time.Now().UnixMilli())

	archives, err := h.repo.QueryArchives(c.Request.Context(), deviceID, from, to)
	if err != nil {
		h.log.Error("query archives failed", "deviceId", deviceID, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"archives": archives,
	})
}

// parseQueryInt64 parses an int64 query parameter with a default value.
func parseQueryInt64(c *gin.Context, key string, defaultVal int64) int64 {
	v := c.Query(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return defaultVal
	}
	return n
}

// parseQueryInt parses an int query parameter with a default value.
func parseQueryInt(c *gin.Context, key string, defaultVal int) int {
	v := c.Query(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return defaultVal
	}
	return n
}
