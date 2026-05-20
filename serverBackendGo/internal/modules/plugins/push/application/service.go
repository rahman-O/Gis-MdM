package application

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/domain"
	pluginpostgres "github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/adapter/persistence/postgres"
	notifpostgres "github.com/gis-mdm/server-backend-go/internal/modules/notifications/adapter/persistence/postgres"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

var ErrPermissionDenied = errors.New("error.permission.denied")

type Service struct {
	repo  *pluginpostgres.MessageRepository
	queue *notifpostgres.QueueRepository
}

func NewService(repo *pluginpostgres.MessageRepository, queue *notifpostgres.QueueRepository) *Service {
	return &Service{repo: repo, queue: queue}
}

func (s *Service) Search(ctx context.Context, p *platformauth.Principal, f domain.PushMessageFilter) (domain.PaginatedMessages, error) {
	items, total, err := s.repo.Search(ctx, int64(p.CustomerID), f)
	if err != nil {
		return domain.PaginatedMessages{}, err
	}
	return domain.PaginatedMessages{Items: items, Total: total}, nil
}

func (s *Service) Send(ctx context.Context, p *platformauth.Principal, req domain.PushSendRequest) error {
	if !p.CanPluginPushSend() {
		return ErrPermissionDenied
	}
	deviceID, err := s.repo.DeviceIDByNumber(ctx, int64(p.CustomerID), req.DeviceNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("device not found")
		}
		return err
	}
	if err := s.queue.Enqueue(ctx, deviceID, req.MessageType, req.Payload); err != nil {
		return err
	}
	return s.repo.InsertHistory(ctx, int64(p.CustomerID), deviceID, req.MessageType, req.Payload)
}

func (s *Service) Delete(ctx context.Context, p *platformauth.Principal, id int64) error {
	if !p.CanPluginPushDelete() {
		return ErrPermissionDenied
	}
	return s.repo.Delete(ctx, int64(p.CustomerID), id)
}

func (s *Service) Purge(ctx context.Context, p *platformauth.Principal, days int) (int64, error) {
	if !p.CanPluginPushDelete() {
		return 0, ErrPermissionDenied
	}
	return s.repo.Purge(ctx, int64(p.CustomerID), days)
}
