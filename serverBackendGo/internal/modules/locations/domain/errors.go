package domain

import "errors"

var (
	ErrBatchSizeInvalid   = errors.New("batch size must be between 1 and 500")
	ErrInvalidCoordinates = errors.New("invalid coordinates: latitude must be -90 to 90, longitude -180 to 180")
	ErrDeviceNotFound     = errors.New("device not found")
	ErrRateLimited        = errors.New("rate limit exceeded")
	ErrDuplicate          = errors.New("duplicate location record")
)
