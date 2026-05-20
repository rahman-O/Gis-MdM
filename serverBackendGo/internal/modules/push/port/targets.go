package port

import "context"

// TargetResolver returns device IDs for push delivery in tenant scope.
type TargetResolver interface {
	DeviceIDsByNumbers(ctx context.Context, customerID, userID int64, allDevices bool, numbers []string) ([]int64, error)
	DeviceIDsByGroupNames(ctx context.Context, customerID, userID int64, allDevices bool, groups []string) ([]int64, error)
	AllDeviceIDs(ctx context.Context, customerID, userID int64, allDevices bool) ([]int64, error)
}

// MessageQueue enqueues agent messages (implemented by notifications postgres repo).
type MessageQueue interface {
	Enqueue(ctx context.Context, deviceID int64, messageType, payload string) error
}
