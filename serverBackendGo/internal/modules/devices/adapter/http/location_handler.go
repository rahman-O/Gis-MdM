package http

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gis-mdm/server-backend-go/internal/platform/httpx/response"
)

// LocationPoint represents a single GPS location reading from the device agent.
type LocationPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
	Speed     float64 `json:"speed"`
	Timestamp int64   `json:"timestamp"`
}

// LocationHandler serves device location endpoints (public + private).
type LocationHandler struct {
	db *sql.DB
}

// NewLocationHandler creates a new LocationHandler.
func NewLocationHandler(db *sql.DB) *LocationHandler {
	return &LocationHandler{db: db}
}

// RegisterPublic registers public endpoints (agent calls these).
func (h *LocationHandler) RegisterPublic(g *gin.RouterGroup) {
	g.POST("/device-locations/:deviceId", h.PostLocations)
}

// RegisterPrivate registers private endpoints (admin panel calls these).
func (h *LocationHandler) RegisterPrivate(g *gin.RouterGroup) {
	g.GET("/devices/:id/locations", h.GetLocations)
}

// PostLocations handles batch location uploads from the device agent.
// POST /public/device-locations/:deviceId
// Body: [{"latitude": 33.3, "longitude": 44.4, "accuracy": 10, "speed": 0, "timestamp": 1716595200000}]
func (h *LocationHandler) PostLocations(c *gin.Context) {
	deviceNumber := c.Param("deviceId")

	var points []LocationPoint
	if err := c.ShouldBindJSON(&points); err != nil {
		response.ErrorEnvelope(c, "invalid body")
		return
	}

	// Find device ID by number
	var deviceID int
	err := h.db.QueryRowContext(c.Request.Context(),
		"SELECT id FROM devices WHERE lower(number) = lower($1)", deviceNumber).Scan(&deviceID)
	if err != nil {
		response.ObjectNotFound(c)
		return
	}

	// Insert all points
	for _, p := range points {
		_, _ = h.db.ExecContext(c.Request.Context(), `
			INSERT INTO device_locations (device_id, latitude, longitude, accuracy, speed, timestamp)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			deviceID, p.Latitude, p.Longitude, p.Accuracy, p.Speed, p.Timestamp)
	}

	// Update device.info with latest location
	if len(points) > 0 {
		latest := points[len(points)-1]
		locJSON := fmt.Sprintf(`{"lat":%f,"lon":%f,"accuracy":%f,"ts":%d}`,
			latest.Latitude, latest.Longitude, latest.Accuracy, latest.Timestamp)
		_, _ = h.db.ExecContext(c.Request.Context(), `
			UPDATE devices SET info = jsonb_set(COALESCE(info::jsonb, '{}'::jsonb), '{location}', $1::jsonb)
			WHERE id = $2`, locJSON, deviceID)
	}

	response.OK(c, gin.H{"saved": len(points)})
}

// GetLocations returns location history for a device.
// GET /private/devices/:id/locations?from=&to=&limit=
func (h *LocationHandler) GetLocations(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	from, _ := strconv.ParseInt(c.Query("from"), 10, 64)
	to, _ := strconv.ParseInt(c.Query("to"), 10, 64)
	limit := 500
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 && l <= 5000 {
		limit = l
	}

	query := `SELECT latitude, longitude, accuracy, speed, timestamp FROM device_locations WHERE device_id = $1`
	args := []any{id}
	argIdx := 2

	if from > 0 {
		query += fmt.Sprintf(` AND timestamp >= $%d`, argIdx)
		args = append(args, from)
		argIdx++
	}
	if to > 0 {
		query += fmt.Sprintf(` AND timestamp <= $%d`, argIdx)
		args = append(args, to)
		argIdx++
	}
	query += fmt.Sprintf(` ORDER BY timestamp DESC LIMIT $%d`, argIdx)
	args = append(args, limit)

	rows, err := h.db.QueryContext(c.Request.Context(), query, args...)
	if err != nil {
		response.ErrorEnvelope(c, "error.internal.server")
		return
	}
	defer rows.Close()

	points := make([]LocationPoint, 0)
	for rows.Next() {
		var p LocationPoint
		if err := rows.Scan(&p.Latitude, &p.Longitude, &p.Accuracy, &p.Speed, &p.Timestamp); err != nil {
			continue
		}
		points = append(points, p)
	}

	response.OK(c, points)
}
