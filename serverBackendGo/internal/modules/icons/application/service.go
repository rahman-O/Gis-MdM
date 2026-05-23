package application

import (
	"context"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/icons/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/icons/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

var ErrPermissionDenied = errors.New("error.permission.denied")

type Service struct {
	repo port.IconRepository
}

func NewService(repo port.IconRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Search(ctx context.Context, p *platformauth.Principal, filter string) ([]domain.Icon, error) {
	if p == nil {
		return nil, ErrPermissionDenied
	}
	return s.repo.List(ctx, p.CustomerID, filter)
}

func (s *Service) Save(ctx context.Context, p *platformauth.Principal, icon domain.Icon) (*domain.Icon, error) {
	if p == nil {
		return nil, ErrPermissionDenied
	}
	icon.CustomerID = p.CustomerID
	return s.repo.Save(ctx, icon)
}

func (s *Service) Delete(ctx context.Context, p *platformauth.Principal, id int) error {
	if p == nil || !p.HasPermission(platformauth.PermSettings) {
		return ErrPermissionDenied
	}
	return s.repo.Delete(ctx, p.CustomerID, id)
}
