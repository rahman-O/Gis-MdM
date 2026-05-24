package application

import (
	"context"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

func (s *Service) Impact(ctx context.Context, p *platformauth.Principal, id int) (*domain.EnrollmentDeleteImpact, error) {
	if err := s.requirePerm(p); err != nil {
		return nil, err
	}
	view, err := s.repo.GetViewByID(ctx, customerID(p), id)
	if err != nil {
		return nil, err
	}
	if view == nil {
		return nil, ErrRouteNotFound
	}
	return s.repo.DeleteImpact(ctx, customerID(p), id)
}

func (s *Service) Delete(ctx context.Context, p *platformauth.Principal, id int) error {
	if err := s.requirePerm(p); err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, customerID(p), id); err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return ErrRouteNotFound
		}
		return err
	}
	return nil
}

func (s *Service) ListTreeNodeOptions(ctx context.Context, p *platformauth.Principal) ([]domain.TreeNodeOption, error) {
	if err := s.requirePerm(p); err != nil {
		return nil, err
	}
	items, err := s.repo.ListTreeNodeOptions(ctx, customerID(p), s.heavyThreshold)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.TreeNodeOption{}
	}
	return items, nil
}

func (s *Service) ListBootstrapApps(ctx context.Context, p *platformauth.Principal) ([]domain.BootstrapAppOption, error) {
	if err := s.requirePerm(p); err != nil {
		return nil, err
	}
	items, err := s.repo.ListBootstrapApps(ctx, customerID(p))
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.BootstrapAppOption{}
	}
	return items, nil
}
