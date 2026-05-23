package application

import (
	"context"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/hints/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/hints/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// Service implements hints use cases.
type Service struct {
	repo port.Repository
}

func NewService(repo port.Repository) *Service {
	return &Service{repo: repo}
}

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrEmptyHintKey    = domain.ErrEmptyHintKey
)

func (s *Service) GetHistory(ctx context.Context, p *platformauth.Principal) ([]string, error) {
	if p == nil {
		return nil, ErrUnauthenticated
	}
	keys, err := s.repo.GetHistory(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	if keys == nil {
		keys = []string{}
	}
	return keys, nil
}

func (s *Service) MarkShown(ctx context.Context, p *platformauth.Principal, rawKey string) error {
	if p == nil {
		return ErrUnauthenticated
	}
	key, err := domain.ValidateHintKey(rawKey)
	if err != nil {
		return err
	}
	return s.repo.MarkShown(ctx, p.ID, string(key))
}

func (s *Service) Enable(ctx context.Context, p *platformauth.Principal) error {
	if p == nil {
		return ErrUnauthenticated
	}
	return s.repo.Enable(ctx, p.ID)
}

func (s *Service) Disable(ctx context.Context, p *platformauth.Principal) error {
	if p == nil {
		return ErrUnauthenticated
	}
	return s.repo.Disable(ctx, p.ID)
}
