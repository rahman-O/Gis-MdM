package application

import (
	"context"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/platform/port"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/platform/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/shared/status"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

var ErrPermissionDenied = errors.New("error.permission.denied")

type Service struct {
	repo  port.Repository
	cache *status.Cache
}

func NewService(repo port.Repository, cache *status.Cache) *Service {
	return &Service{repo: repo, cache: cache}
}

func (s *Service) Active(ctx context.Context) ([]domain.Plugin, error) {
	return s.repo.FindActive(ctx)
}

func (s *Service) Available(ctx context.Context, p *platformauth.Principal) ([]domain.Plugin, error) {
	if p == nil {
		return nil, ErrPermissionDenied
	}
	return s.repo.FindAvailableByCustomer(ctx, int64(p.CustomerID))
}

func (s *Service) Registered(ctx context.Context) ([]domain.Plugin, error) {
	return s.repo.FindRegistered(ctx)
}

func (s *Service) SaveDisabled(ctx context.Context, p *platformauth.Principal, pluginIDs []int64) error {
	if p == nil || !p.CanManagePluginsCustomer() {
		return ErrPermissionDenied
	}
	if err := s.repo.SaveDisabled(ctx, int64(p.CustomerID), pluginIDs); err != nil {
		return err
	}
	s.cache.SetDisabled(int64(p.CustomerID), pluginIDs)
	return nil
}
