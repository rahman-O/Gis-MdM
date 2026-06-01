package application

import (
	"context"
	"log/slog"
	"math"

	"github.com/gis-mdm/server-backend-go/internal/modules/locations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/locations/port"
)

const (
	coordTolerance    = 0.00001 // ~1.1 meters
	timeTolerance     = 10_000  // 10 seconds in ms
	distanceTolerance = 50.0    // 50 meters
	earthRadiusKm     = 6371.0
)

// HaversineDistance calculates the distance in meters between two GPS coordinates.
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c * 1000 // meters
}

// DuplicateDetector checks if an incoming location record is a duplicate.
type DuplicateDetector struct {
	repo port.Repository
	log  *slog.Logger
}

func NewDuplicateDetector(repo port.Repository, log *slog.Logger) *DuplicateDetector {
	return &DuplicateDetector{repo: repo, log: log}
}

// IsDuplicate returns true if the record is too close to the last stored record.
func (d *DuplicateDetector) IsDuplicate(ctx context.Context, deviceID int, incoming domain.LocationRecord) (bool, string) {
	last, err := d.repo.GetLastRecord(ctx, deviceID)
	if err != nil {
		d.log.Warn("dedup: failed to get last record, allowing through", "deviceId", deviceID, "err", err)
		return false, "" // fail-open
	}
	if last == nil {
		return false, "" // first record
	}

	// Check coordinate tolerance
	if math.Abs(incoming.Latitude-last.Latitude) < coordTolerance &&
		math.Abs(incoming.Longitude-last.Longitude) < coordTolerance {
		return true, "coordinate_tolerance"
	}

	// Check timestamp proximity + distance
	if math.Abs(float64(incoming.Timestamp-last.Timestamp)) < timeTolerance {
		dist := HaversineDistance(last.Latitude, last.Longitude, incoming.Latitude, incoming.Longitude)
		if dist < distanceTolerance {
			return true, "timestamp_proximity"
		}
	}

	return false, ""
}
