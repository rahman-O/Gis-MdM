package application

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/notifications/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/notifications/port"
)

var (
	ErrDeviceNotFound = errors.New("error.notfound.device")
)

type Service struct {
	devices port.DeviceLookup
	queue   port.MessageQueue
}

func NewService(devices port.DeviceLookup, queue port.MessageQueue) *Service {
	return &Service{devices: devices, queue: queue}
}

func (s *Service) GetPending(ctx context.Context, deviceNumber string) ([]domain.PlainPushMessage, error) {
	deviceID, err := s.devices.DeviceIDByNumber(ctx, deviceNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}
	msgs, err := s.queue.ListPendingForDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if len(msgs) > 0 {
		ids := make([]int64, len(msgs))
		for i := range msgs {
			ids[i] = msgs[i].ID
		}
		_ = s.queue.MarkDelivered(ctx, ids)
	}
	return msgs, nil
}

func (s *Service) PollPending(ctx context.Context, deviceNumber string) ([]domain.PlainPushMessage, error) {
	return s.GetPending(ctx, deviceNumber)
}
