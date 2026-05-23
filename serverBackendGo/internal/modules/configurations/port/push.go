package port

// PushNotifier is a no-op until push module is migrated.
type PushNotifier interface {
	NotifyConfigurationChanged(configurationID int) error
}

// NoopPushNotifier does nothing.
type NoopPushNotifier struct{}

func (NoopPushNotifier) NotifyConfigurationChanged(int) error { return nil }
