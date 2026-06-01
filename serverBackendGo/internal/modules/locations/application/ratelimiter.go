package application

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// RateLimitResult contains the result of a rate limit check.
type RateLimitResult struct {
	Allowed      bool
	CurrentCount int
	Limit        int
	RetryAfterMs int64
}

// LocationRateLimiter enforces per-device write rate limits using a sliding window.
type LocationRateLimiter struct {
	mu        sync.Mutex
	windows   map[string][]int64
	maxWrites int
	windowMs  int64
	log       *slog.Logger
}

func NewLocationRateLimiter(maxWrites int, windowSec int, log *slog.Logger) *LocationRateLimiter {
	return &LocationRateLimiter{
		windows:   make(map[string][]int64),
		maxWrites: maxWrites,
		windowMs:  int64(windowSec) * 1000,
		log:       log,
	}
}

// CheckLimit checks if a device is within its rate limit.
func (r *LocationRateLimiter) CheckLimit(deviceID string) RateLimitResult {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UnixMilli()
	cutoff := now - r.windowMs

	// Remove expired entries
	timestamps := r.windows[deviceID]
	valid := timestamps[:0]
	for _, ts := range timestamps {
		if ts > cutoff {
			valid = append(valid, ts)
		}
	}
	r.windows[deviceID] = valid

	if len(valid) >= r.maxWrites {
		r.log.Warn("rate limit exceeded", "deviceId", deviceID, "count", len(valid), "limit", r.maxWrites)
		oldest := valid[0]
		retryAfter := oldest + r.windowMs - now
		return RateLimitResult{Allowed: false, CurrentCount: len(valid), Limit: r.maxWrites, RetryAfterMs: retryAfter}
	}

	r.windows[deviceID] = append(valid, now)
	return RateLimitResult{Allowed: true, CurrentCount: len(valid) + 1, Limit: r.maxWrites}
}

// StartCleanup runs a periodic goroutine to remove stale device entries.
func (r *LocationRateLimiter) StartCleanup(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.cleanup()
		}
	}
}

func (r *LocationRateLimiter) cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UnixMilli()
	cutoff := now - r.windowMs
	for deviceID, timestamps := range r.windows {
		valid := timestamps[:0]
		for _, ts := range timestamps {
			if ts > cutoff {
				valid = append(valid, ts)
			}
		}
		if len(valid) == 0 {
			delete(r.windows, deviceID)
		} else {
			r.windows[deviceID] = valid
		}
	}
}
