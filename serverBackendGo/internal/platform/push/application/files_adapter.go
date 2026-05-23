package application

// FilesPushAdapter adapts Notifier to files/port.PushNotifier (void return).
type FilesPushAdapter struct {
	*Notifier
}

func (a FilesPushAdapter) NotifyConfigurationUpdate(configurationID int) {
	if a.Notifier != nil {
		_ = a.Notifier.NotifyConfigurationUpdate(configurationID)
	}
}
