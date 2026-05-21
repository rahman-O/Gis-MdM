package application

import (
	"context"
	"log/slog"

	pushport "github.com/gis-mdm/server-backend-go/internal/platform/push/port"
)

const (
	TypeConfigUpdated    = "configUpdated"
	TypeAppConfigUpdated = "appConfigUpdated"
)

// Notifier enqueues push messages for devices (Java PushService polling path).
type Notifier struct {
	queue  pushport.MessageQueue
	lookup pushport.DeviceLookup
	log    *slog.Logger
}

func NewNotifier(queue pushport.MessageQueue, lookup pushport.DeviceLookup, log *slog.Logger) *Notifier {
	return &Notifier{queue: queue, lookup: lookup, log: log}
}

// NotifyConfigurationChanged implements configurations/port.PushNotifier.
func (n *Notifier) NotifyConfigurationChanged(configurationID int) error {
	return n.notifyConfiguration(context.Background(), configurationID)
}

// NotifyConfigurationUpdate implements files/port.PushNotifier.
func (n *Notifier) NotifyConfigurationUpdate(configurationID int) error {
	return n.notifyConfiguration(context.Background(), configurationID)
}

func (n *Notifier) notifyConfiguration(ctx context.Context, configurationID int) error {
	if n == nil || n.queue == nil || n.lookup == nil {
		return nil
	}
	ids, err := n.lookup.DeviceIDsByConfiguration(ctx, configurationID)
	if err != nil {
		n.logWarn("configuration push lookup failed", configurationID, err)
		return nil
	}
	for _, id := range ids {
		if err := n.queue.Enqueue(ctx, id, TypeConfigUpdated, ""); err != nil {
			n.logWarn("configuration push enqueue failed", configurationID, err)
		}
	}
	return nil
}

// NotifyAppSettings implements devices/port.PushNotifier.
func (n *Notifier) NotifyAppSettings(ctx context.Context, deviceID int) error {
	if n == nil || n.queue == nil {
		return nil
	}
	if err := n.queue.Enqueue(ctx, int64(deviceID), TypeAppConfigUpdated, ""); err != nil {
		n.logWarn("app settings push enqueue failed", deviceID, err)
	}
	return nil
}

func (n *Notifier) logWarn(msg string, ref int, err error) {
	if n.log != nil {
		n.log.Warn(msg, slog.Int("ref", ref), slog.String("error", err.Error()))
	}
}
