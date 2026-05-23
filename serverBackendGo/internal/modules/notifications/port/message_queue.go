package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/notifications/domain"
)

// MessageQueue persists agent push messages (pushmessages + pendingpushes).
type MessageQueue interface {
	Enqueue(ctx context.Context, deviceID int64, messageType, payload string) error
	ListPendingForDevice(ctx context.Context, deviceID int64) ([]domain.PlainPushMessage, error)
	MarkDelivered(ctx context.Context, messageIDs []int64) error
}

// DeviceLookup resolves device id by number (current or old).
type DeviceLookup interface {
	DeviceIDByNumber(ctx context.Context, deviceNumber string) (int64, error)
}
