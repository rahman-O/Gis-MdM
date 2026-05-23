package application

import (
	"context"
	"errors"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/groups/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/groups/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// Service implements group use cases.
type Service struct {
	repo port.GroupRepository
}

func NewService(repo port.GroupRepository) *Service {
	return &Service{repo: repo}
}

var (
	ErrPermissionDenied = errors.New("error.permission.denied")
	ErrDuplicateGroup   = errors.New("error.duplicate.group")
	ErrNotEmptyGroup    = errors.New("error.notempty.group")
)

func customerID(p *platformauth.Principal) int {
	if p == nil {
		return 0
	}
	return p.CustomerID
}

func (s *Service) List(ctx context.Context, p *platformauth.Principal) ([]domain.Group, error) {
	if p == nil {
		return nil, ErrPermissionDenied
	}
	return s.repo.ListByCustomer(ctx, customerID(p))
}

func (s *Service) SearchByValue(ctx context.Context, p *platformauth.Principal, value string) ([]domain.Group, error) {
	if p == nil {
		return nil, ErrPermissionDenied
	}
	return s.repo.ListByValue(ctx, customerID(p), value)
}

func (s *Service) Autocomplete(ctx context.Context, p *platformauth.Principal, filter string) ([]domain.LookupItem, error) {
	groups, err := s.SearchByValue(ctx, p, filter)
	if err != nil {
		return nil, err
	}
	out := make([]domain.LookupItem, 0, len(groups))
	for _, g := range groups {
		name := g.Name
		out = append(out, domain.LookupItem{ID: g.ID, Name: &name})
	}
	return out, nil
}

func (s *Service) Save(ctx context.Context, p *platformauth.Principal, g domain.Group) error {
	if p == nil || !p.HasPermission(platformauth.PermSettings) {
		return ErrPermissionDenied
	}
	g.Name = strings.TrimSpace(g.Name)
	if g.Name == "" {
		return ErrPermissionDenied
	}
	cid := customerID(p)
	existing, err := s.repo.GetByName(ctx, cid, g.Name)
	if err != nil {
		return err
	}
	if existing != nil && existing.ID != g.ID {
		return ErrDuplicateGroup
	}
	if g.ID == 0 {
		id, err := s.repo.Insert(ctx, cid, g.Name)
		if err != nil {
			return err
		}
		all, err := s.repo.UserHasAllDevices(ctx, p.ID)
		if err != nil {
			return err
		}
		if !all {
			return s.repo.GrantCreatorAccess(ctx, p.ID, id)
		}
		return nil
	}
	return s.repo.Update(ctx, cid, g)
}

func (s *Service) Delete(ctx context.Context, p *platformauth.Principal, id int) error {
	if p == nil || !p.HasPermission(platformauth.PermSettings) {
		return ErrPermissionDenied
	}
	n, err := s.repo.CountDevicesInGroup(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrNotEmptyGroup
	}
	return s.repo.Delete(ctx, customerID(p), id)
}
