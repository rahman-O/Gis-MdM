package application

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/notifications/port"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/messaging/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/messaging/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/shared/targets"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

var ErrPermissionDenied = errors.New("error.permission.denied")

type Service struct {
	repo    *postgres.MessageRepository
	queue   port.MessageQueue
	targets targets.Resolver
}

func NewService(repo *postgres.MessageRepository, queue port.MessageQueue, tr targets.Resolver) *Service {
	return &Service{repo: repo, queue: queue, targets: tr}
}

func (s *Service) Search(ctx context.Context, p *platformauth.Principal, f domain.MessageFilter) (domain.PaginatedMessages, error) {
	if p == nil {
		return domain.PaginatedMessages{}, ErrPermissionDenied
	}
	items, total, err := s.repo.Search(ctx, int64(p.CustomerID), f)
	if err != nil {
		return domain.PaginatedMessages{}, err
	}
	return domain.PaginatedMessages{Items: items, Total: total}, nil
}

func (s *Service) Send(ctx context.Context, p *platformauth.Principal, req domain.SendRequest) error {
	if p == nil || !p.CanPluginMessagingSend() {
		return ErrPermissionDenied
	}
	cid := int64(p.CustomerID)
	uid := p.ID
	var deviceIDs []int64
	var err error
	scope := strings.ToLower(strings.TrimSpace(req.Scope))
	switch scope {
	case "device", "":
		if req.DeviceNumber == "" {
			return errors.New("error.params.missing")
		}
		deviceIDs, err = s.targets.DeviceIDsByNumbers(ctx, cid, uid, []string{req.DeviceNumber})
	case "group":
		if req.GroupID == 0 {
			return errors.New("error.params.missing")
		}
		deviceIDs, err = s.targets.DeviceIDsByGroupID(ctx, cid, uid, req.GroupID)
	case "configuration":
		if req.ConfigurationID == 0 {
			return errors.New("error.params.missing")
		}
		deviceIDs, err = s.targets.DeviceIDsByConfigurationID(ctx, cid, uid, req.ConfigurationID)
	default:
		return errors.New("error.params.missing")
	}
	if err != nil {
		return err
	}
	if len(deviceIDs) == 0 {
		return errors.New("device not found")
	}
	for _, did := range deviceIDs {
		id, err := s.repo.Insert(ctx, cid, did, req.Message)
		if err != nil {
			return err
		}
		payload := postgres.QueuePayload(id, req.Message)
		if err := s.queue.Enqueue(ctx, did, "textMessage", payload); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, p *platformauth.Principal, id int64) error {
	if p == nil || !p.CanPluginMessagingDelete() {
		return ErrPermissionDenied
	}
	return s.repo.Delete(ctx, int64(p.CustomerID), id)
}

func (s *Service) Purge(ctx context.Context, p *platformauth.Principal, days int) error {
	if p == nil || !p.CanPluginMessagingDelete() {
		return ErrPermissionDenied
	}
	return s.repo.Purge(ctx, int64(p.CustomerID), days)
}

func (s *Service) SetStatus(ctx context.Context, id int64, status int) error {
	if status < 0 || status > domain.StatusRead {
		return errors.New("error.params.invalid")
	}
	err := s.repo.UpdateStatus(ctx, id, status)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("error.notfound")
	}
	return err
}
