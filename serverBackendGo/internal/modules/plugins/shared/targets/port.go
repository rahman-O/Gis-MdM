package targets

import "context"

// Resolver resolves device IDs for push/messaging targets.
type Resolver interface {
	DeviceIDsByNumbers(ctx context.Context, customerID, userID int64, numbers []string) ([]int64, error)
	DeviceIDsByGroupID(ctx context.Context, customerID, userID, groupID int64) ([]int64, error)
	DeviceIDsByConfigurationID(ctx context.Context, customerID, userID, configurationID int64) ([]int64, error)
	AllDeviceIDs(ctx context.Context, customerID, userID int64) ([]int64, error)
}
