package application

import (
	"context"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/push/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/push/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

var ErrPermissionDenied = errors.New("error.permission.denied")

type Service struct {
	targets port.TargetResolver
	queue   port.MessageQueue
}

func NewService(targets port.TargetResolver, queue port.MessageQueue) *Service {
	return &Service{targets: targets, queue: queue}
}

func (s *Service) Send(ctx context.Context, p *platformauth.Principal, req domain.PushRequest) error {
	if p == nil || !p.CanUsePushAPI() {
		return ErrPermissionDenied
	}
	ids := map[int64]struct{}{}
	add := func(list []int64) {
		for _, id := range list {
			ids[id] = struct{}{}
		}
	}
	if req.Broadcast {
		list, err := s.targets.AllDeviceIDs(ctx, int64(p.CustomerID), p.ID, true)
		if err != nil {
			return err
		}
		add(list)
	} else {
		if len(req.Groups) > 0 {
			list, err := s.targets.DeviceIDsByGroupNames(ctx, int64(p.CustomerID), p.ID, true, req.Groups)
			if err != nil {
				return err
			}
			add(list)
		}
		if len(req.DeviceNumbers) > 0 {
			list, err := s.targets.DeviceIDsByNumbers(ctx, int64(p.CustomerID), p.ID, true, req.DeviceNumbers)
			if err != nil {
				return err
			}
			add(list)
		}
	}
	for id := range ids {
		if err := s.queue.Enqueue(ctx, id, req.MessageType, req.Payload); err != nil {
			return err
		}
	}
	return nil
}
