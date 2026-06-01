package application

import (
	"context"
	"log/slog"
	"sort"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/locations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/locations/port"
)

const (
	defaultRetentionDays      = 30
	defaultAggregationHours   = 24
	defaultAggregationBatch   = 10000
	msPerHour                 = 3_600_000
)

// AggregationResult holds the outcome of a single aggregation run.
type AggregationResult struct {
	Success          bool
	RecordsArchived  int
	SummariesCreated int
	Error            string
}

// AggregationWorker archives old location records into hourly summaries.
type AggregationWorker struct {
	repo          port.Repository
	log           *slog.Logger
	retentionDays int
	intervalHours int
	batchSize     int
}

// NewAggregationWorker creates a new AggregationWorker.
func NewAggregationWorker(repo port.Repository, log *slog.Logger, retentionDays, intervalHours int) *AggregationWorker {
	if retentionDays <= 0 {
		retentionDays = defaultRetentionDays
	}
	if intervalHours <= 0 {
		intervalHours = defaultAggregationHours
	}
	return &AggregationWorker{
		repo:          repo,
		log:           log,
		retentionDays: retentionDays,
		intervalHours: intervalHours,
		batchSize:     defaultAggregationBatch,
	}
}

// Start runs the aggregation loop until the context is cancelled.
func (w *AggregationWorker) Start(ctx context.Context) {
	w.log.Info("aggregation worker started",
		"retentionDays", w.retentionDays,
		"intervalHours", w.intervalHours,
	)

	// Run once on startup after a short delay
	select {
	case <-ctx.Done():
		return
	case <-time.After(30 * time.Second):
		w.RunOnce(ctx)
	}

	ticker := time.NewTicker(time.Duration(w.intervalHours) * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.log.Info("aggregation worker stopped")
			return
		case <-ticker.C:
			w.RunOnce(ctx)
		}
	}
}

// RunOnce performs a single aggregation cycle.
func (w *AggregationWorker) RunOnce(ctx context.Context) AggregationResult {
	cutoff := time.Now().Add(-time.Duration(w.retentionDays) * 24 * time.Hour).UnixMilli()

	records, err := w.repo.GetRecordsOlderThan(ctx, cutoff, w.batchSize)
	if err != nil {
		w.log.Error("aggregation: failed to query old records", "err", err)
		return AggregationResult{Error: err.Error()}
	}

	if len(records) == 0 {
		w.log.Debug("aggregation: no records to archive")
		return AggregationResult{Success: true}
	}

	// Group records by device + hour
	type groupKey struct {
		DeviceID  int
		HourStart int64
	}
	groups := make(map[groupKey][]domain.LocationRecord)
	for _, rec := range records {
		hourStart := (rec.Timestamp / msPerHour) * msPerHour
		key := groupKey{DeviceID: rec.DeviceID, HourStart: hourStart}
		groups[key] = append(groups[key], rec)
	}

	// Build archives
	archives := make([]domain.LocationArchive, 0, len(groups))
	for key, recs := range groups {
		archive := w.buildArchive(key.DeviceID, key.HourStart, recs)
		archives = append(archives, archive)
	}

	// Insert archives
	if err := w.repo.InsertArchives(ctx, archives); err != nil {
		w.log.Error("aggregation: failed to insert archives", "err", err, "count", len(archives))
		return AggregationResult{Error: err.Error()}
	}

	// Delete originals
	ids := make([]int64, 0, len(records))
	for _, rec := range records {
		ids = append(ids, rec.ID)
	}
	if err := w.repo.DeleteByIDs(ctx, ids); err != nil {
		w.log.Error("aggregation: failed to delete originals", "err", err, "count", len(ids))
		return AggregationResult{
			Success:          true, // archives were inserted
			RecordsArchived:  len(records),
			SummariesCreated: len(archives),
			Error:            "delete_failed: " + err.Error(),
		}
	}

	w.log.Info("aggregation complete",
		"recordsArchived", len(records),
		"summariesCreated", len(archives),
	)

	return AggregationResult{
		Success:          true,
		RecordsArchived:  len(records),
		SummariesCreated: len(archives),
	}
}

// buildArchive computes an hourly summary from a set of records.
func (w *AggregationWorker) buildArchive(deviceID int, hourStart int64, records []domain.LocationRecord) domain.LocationArchive {
	// Sort by timestamp
	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp < records[j].Timestamp
	})

	first := records[0]
	last := records[len(records)-1]

	// Compute total distance traveled
	var totalDistance float64
	for i := 1; i < len(records); i++ {
		totalDistance += HaversineDistance(
			records[i-1].Latitude, records[i-1].Longitude,
			records[i].Latitude, records[i].Longitude,
		)
	}

	return domain.LocationArchive{
		DeviceID:         deviceID,
		HourStart:        hourStart,
		StartLatitude:    first.Latitude,
		StartLongitude:   first.Longitude,
		EndLatitude:      last.Latitude,
		EndLongitude:     last.Longitude,
		DistanceTraveled: totalDistance,
		PointCount:       len(records),
	}
}
