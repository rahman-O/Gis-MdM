package port

import "context"

// MessageQueue enqueues agent notifications (implemented by notifications postgres repo).
type MessageQueue interface {
	Enqueue(ctx context.Context, deviceID int64, messageType, payload string) error
}

// DeviceLookup resolves devices for configuration-scoped push.
type DeviceLookup interface {
	DeviceIDsByConfiguration(ctx context.Context, configurationID int) ([]int64, error)
}
