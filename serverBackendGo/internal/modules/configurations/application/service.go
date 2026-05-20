package application

import (
	"context"
	"errors"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/configurations/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// Service implements configuration use cases.
type Service struct {
	repo port.ConfigRepository
	push port.PushNotifier
}

func NewService(repo port.ConfigRepository, push port.PushNotifier) *Service {
	return &Service{repo: repo, push: push}
}

var (
	ErrPermissionDenied      = errors.New("error.permission.denied")
	ErrDuplicateConfiguration  = errors.New("error.duplicate.configuration")
	ErrNotEmptyConfiguration   = errors.New("error.notempty.configuration")
	ErrConfigurationNotFound   = errors.New("error.notfound.configuration")
)

func customerID(p *platformauth.Principal) int {
	if p == nil {
		return 0
	}
	return p.CustomerID
}

func (s *Service) ListNames(ctx context.Context, p *platformauth.Principal) ([]domain.LookupItem, error) {
	if p == nil {
		return nil, ErrPermissionDenied
	}
	return s.repo.ListByCustomer(ctx, customerID(p))
}

func (s *Service) Search(ctx context.Context, p *platformauth.Principal) ([]domain.Configuration, error) {
	if p == nil || !p.CanManageConfigurations() {
		return nil, ErrPermissionDenied
	}
	return s.repo.Search(ctx, customerID(p))
}

func (s *Service) SearchByValue(ctx context.Context, p *platformauth.Principal, value string) ([]domain.Configuration, error) {
	if p == nil || !p.CanManageConfigurations() {
		return nil, ErrPermissionDenied
	}
	return s.repo.SearchByValue(ctx, customerID(p), value)
}

func (s *Service) Autocomplete(ctx context.Context, p *platformauth.Principal, filter string) ([]domain.LookupItem, error) {
	cfgs, err := s.SearchByValue(ctx, p, filter)
	if err != nil {
		return nil, err
	}
	out := make([]domain.LookupItem, 0, len(cfgs))
	for _, c := range cfgs {
		if c.Name == nil {
			continue
		}
		id := 0
		if c.ID != nil {
			id = *c.ID
		}
		name := *c.Name
		out = append(out, domain.LookupItem{ID: id, Name: &name})
	}
	return out, nil
}

func (s *Service) GetByID(ctx context.Context, p *platformauth.Principal, id int) (*domain.Configuration, error) {
	if p == nil || !p.CanManageConfigurations() {
		return nil, ErrPermissionDenied
	}
	cfg, err := s.repo.GetByID(ctx, customerID(p), id)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, ErrConfigurationNotFound
	}
	return cfg, nil
}

func (s *Service) Save(ctx context.Context, p *platformauth.Principal, cfg domain.Configuration) (*domain.Configuration, error) {
	if p == nil || !p.CanManageConfigurations() {
		return nil, ErrPermissionDenied
	}
	if cfg.Name != nil {
		n := strings.TrimSpace(*cfg.Name)
		cfg.Name = &n
		if n == "" {
			return nil, ErrPermissionDenied
		}
	}
	cid := customerID(p)
	existing, err := s.repo.GetByName(ctx, cid, derefStr(cfg.Name))
	if err != nil {
		return nil, err
	}
	id := 0
	if cfg.ID != nil {
		id = *cfg.ID
	}
	if existing != nil && existing.ID != nil && *existing.ID != id {
		return nil, ErrDuplicateConfiguration
	}
	var savedID int
	if id == 0 {
		savedID, err = s.repo.Insert(ctx, cid, cfg)
	} else {
		err = s.repo.Update(ctx, cid, cfg)
		savedID = id
	}
	if err != nil {
		return nil, err
	}
	_ = s.push.NotifyConfigurationChanged(savedID)
	return s.repo.GetByID(ctx, cid, savedID)
}

func (s *Service) Delete(ctx context.Context, p *platformauth.Principal, id int) error {
	if p == nil || !p.CanManageConfigurations() {
		return ErrPermissionDenied
	}
	n, err := s.repo.CountDevicesUsing(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrNotEmptyConfiguration
	}
	return s.repo.Delete(ctx, customerID(p), id)
}

func (s *Service) Copy(ctx context.Context, p *platformauth.Principal, req domain.CopyRequest) (int, error) {
	if p == nil || !p.CanManageConfigurations() {
		return 0, ErrPermissionDenied
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return 0, ErrPermissionDenied
	}
	existing, err := s.repo.GetByName(ctx, customerID(p), req.Name)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, ErrDuplicateConfiguration
	}
	return s.repo.Copy(ctx, customerID(p), req)
}

func (s *Service) ListAllApplicationsForPicker(ctx context.Context, p *platformauth.Principal) ([]domain.ConfigurationApplication, error) {
	if p == nil || !p.CanManageConfigurations() {
		return nil, ErrPermissionDenied
	}
	return s.repo.ListAllApplicationsForPicker(ctx, customerID(p))
}

func (s *Service) ListConfigurationApplications(ctx context.Context, p *platformauth.Principal, configurationID int) ([]domain.ConfigurationApplication, error) {
	if p == nil || !p.CanManageConfigurations() {
		return nil, ErrPermissionDenied
	}
	return s.repo.ListConfigurationApplications(ctx, customerID(p), configurationID)
}

func (s *Service) UpgradeApplication(ctx context.Context, p *platformauth.Principal, req domain.UpgradeApplicationRequest) error {
	if p == nil || !p.CanManageConfigurations() {
		return ErrPermissionDenied
	}
	return s.repo.UpgradeApplication(ctx, customerID(p), req)
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
