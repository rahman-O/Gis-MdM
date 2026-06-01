package application

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/locations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/locations/port"
)

// LocationService orchestrates the batch ingestion pipeline.
type LocationService struct {
	repo        port.Repository
	dedup       *DuplicateDetector
	limiter     *LocationRateLimiter
	broadcaster port.Broadcaster
	log         *slog.Logger
	batchMax    int
}

// NewLocationService creates a new LocationService with the given dependencies.
func NewLocationService(
	repo port.Repository,
	dedup *DuplicateDetector,
	limiter *LocationRateLimiter,
	broadcaster port.Broadcaster,
	log *slog.Logger,
	batchMax int,
) *LocationService {
	if batchMax <= 0 {
		batchMax = 500
	}
	return &LocationService{
		repo:        repo,
		dedup:       dedup,
		limiter:     limiter,
		broadcaster: broadcaster,
		log:         log,
		batchMax:    batchMax,
	}
}

// ProcessBatch validates, deduplicates, rate-limits, and stores a batch of location entries.
// It always returns a response (never an error) to prevent agent retries.
func (s *LocationService) ProcessBatch(ctx context.Context, deviceID int, entries []domain.BatchUploadEntry) *domain.BatchUploadResponse {
	resp := &domain.BatchUploadResponse{Status: "OK"}

	if len(entries) == 0 || len(entries) > s.batchMax {
		resp.Rejected = len(entries)
		resp.Reasons = append(resp.Reasons, fmt.Sprintf("batch size must be 1-%d", s.batchMax))
		return resp
	}

	deviceIDStr := strconv.Itoa(deviceID)
	now := time.Now()
	month := now.Format("2006-01")

	var valid []domain.LocationRecord

	for _, entry := range entries {
		// Rate limit check
		result := s.limiter.CheckLimit(deviceIDStr)
		if !result.Allowed {
			resp.Rejected++
			resp.Reasons = append(resp.Reasons, "rate_limited")
			continue
		}

		// Coordinate validation
		if entry.Latitude < -90 || entry.Latitude > 90 || entry.Longitude < -180 || entry.Longitude > 180 {
			resp.Rejected++
			resp.Reasons = append(resp.Reasons, "invalid_coordinates")
			continue
		}

		// Build record
		rec := domain.LocationRecord{
			DeviceID:     deviceID,
			Latitude:     entry.Latitude,
			Longitude:    entry.Longitude,
			Accuracy:     entry.Accuracy,
			Speed:        entry.Speed,
			Altitude:     entry.Altitude,
			BatteryLevel: entry.BatteryLevel,
			NetworkType:  entry.NetworkType,
			TrackingMode: entry.TrackingMode,
			Timestamp:    entry.Timestamp,
			Month:        month,
		}

		// Dedup check
		isDup, reason := s.dedup.IsDuplicate(ctx, deviceID, rec)
		if isDup {
			resp.Rejected++
			resp.Reasons = append(resp.Reasons, "duplicate:"+reason)
			continue
		}

		valid = append(valid, rec)
	}

	if len(valid) > 0 {
		if err := s.repo.InsertBatch(ctx, valid); err != nil {
			s.log.Error("batch insert failed", "deviceId", deviceID, "err", err)
			resp.Rejected += len(valid)
			resp.Reasons = append(resp.Reasons, "insert_failed")
			valid = nil
		}
	}

	resp.Accepted = len(valid)

	// Broadcast to WebSocket subscribers
	if s.broadcaster != nil && len(valid) > 0 {
		s.broadcaster.UpdateDeviceHeartbeat(deviceIDStr)
		last := valid[len(valid)-1]
		s.broadcaster.BroadcastLocation(deviceIDStr, domain.LocationUpdateData{
			Latitude:     last.Latitude,
			Longitude:    last.Longitude,
			Accuracy:     last.Accuracy,
			Speed:        last.Speed,
			Altitude:     last.Altitude,
			BatteryLevel: last.BatteryLevel,
			NetworkType:  last.NetworkType,
			Timestamp:    last.Timestamp,
		})
	}

	return resp
}
