package port

import "context"

// PushNotifier sends device push notifications (no-op in Phase 4).
type PushNotifier interface {
	NotifyAppSettings(ctx context.Context, deviceID int) error
}

// NoopPush is a stub implementation.
type NoopPush struct{}

func (NoopPush) NotifyAppSettings(context.Context, int) error { return nil }
