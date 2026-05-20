package application

import (
	"context"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/applications/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/applications/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// Service implements application use cases.
type Service struct {
	repo port.ApplicationRepository
}

func NewService(repo port.ApplicationRepository) *Service {
	return &Service{repo: repo}
}

var (
	ErrPermissionDenied = errors.New("error.permission.denied")
	ErrAppNotFound      = errors.New("error.notfound.application")
)

func customerID(p *platformauth.Principal) int {
	if p == nil {
		return 0
	}
	return p.CustomerID
}

func (s *Service) Search(ctx context.Context, p *platformauth.Principal) ([]domain.Application, error) {
	if p == nil || !p.CanManageApplications() {
		return nil, ErrPermissionDenied
	}
	return s.repo.Search(ctx, customerID(p))
}

func (s *Service) SearchByValue(ctx context.Context, p *platformauth.Principal, value string) ([]domain.Application, error) {
	if p == nil || !p.CanManageApplications() {
		return nil, ErrPermissionDenied
	}
	return s.repo.SearchByValue(ctx, customerID(p), value)
}

func (s *Service) Autocomplete(ctx context.Context, p *platformauth.Principal, filter string) ([]domain.LookupItem, error) {
	apps, err := s.SearchByValue(ctx, p, filter)
	if err != nil {
		return nil, err
	}
	out := make([]domain.LookupItem, 0, len(apps))
	for _, a := range apps {
		if a.Name == nil || a.ID == nil {
			continue
		}
		name := *a.Name
		out = append(out, domain.LookupItem{ID: *a.ID, Name: &name})
	}
	return out, nil
}

func (s *Service) GetByID(ctx context.Context, p *platformauth.Principal, id int) (*domain.Application, error) {
	if p == nil || !p.CanManageApplications() {
		return nil, ErrPermissionDenied
	}
	app, err := s.repo.GetByID(ctx, customerID(p), id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, ErrAppNotFound
	}
	return app, nil
}

func (s *Service) ListVersions(ctx context.Context, p *platformauth.Principal, applicationID int) ([]domain.ApplicationVersion, error) {
	if p == nil || !p.CanManageApplications() {
		return nil, ErrPermissionDenied
	}
	return s.repo.ListVersions(ctx, customerID(p), applicationID)
}

func (s *Service) SaveAndroid(ctx context.Context, p *platformauth.Principal, app domain.Application) (*domain.Application, error) {
	if p == nil || !p.CanManageApplications() {
		return nil, ErrPermissionDenied
	}
	return s.repo.SaveAndroid(ctx, customerID(p), app)
}

func (s *Service) SaveWeb(ctx context.Context, p *platformauth.Principal, app domain.Application) (*domain.Application, error) {
	if p == nil || !p.CanManageApplications() {
		return nil, ErrPermissionDenied
	}
	return s.repo.SaveWeb(ctx, customerID(p), app)
}

func (s *Service) SaveVersion(ctx context.Context, p *platformauth.Principal, ver domain.ApplicationVersion) (*domain.ApplicationVersion, error) {
	if p == nil || !p.CanManageApplications() {
		return nil, ErrPermissionDenied
	}
	return s.repo.SaveVersion(ctx, customerID(p), ver)
}

func (s *Service) DeleteApp(ctx context.Context, p *platformauth.Principal, id int) error {
	if p == nil || !p.CanManageApplications() {
		return ErrPermissionDenied
	}
	return s.repo.DeleteApp(ctx, customerID(p), id)
}

func (s *Service) DeleteVersion(ctx context.Context, p *platformauth.Principal, versionID int) error {
	if p == nil || !p.CanManageApplications() {
		return ErrPermissionDenied
	}
	return s.repo.DeleteVersion(ctx, customerID(p), versionID)
}

func (s *Service) ValidatePkg(ctx context.Context, p *platformauth.Principal, req domain.ValidatePkgRequest) ([]domain.Application, error) {
	if p == nil || !p.CanManageApplications() {
		return nil, ErrPermissionDenied
	}
	return s.repo.ValidatePkg(ctx, customerID(p), req)
}

func (s *Service) GetAppConfigurations(ctx context.Context, p *platformauth.Principal, applicationID int) ([]domain.ApplicationConfigurationLink, error) {
	if p == nil || !p.CanManageApplications() {
		return nil, ErrPermissionDenied
	}
	return s.repo.GetAppConfigurations(ctx, customerID(p), applicationID)
}

func (s *Service) UpdateAppConfigurations(ctx context.Context, p *platformauth.Principal, req domain.LinkConfigurationsToAppRequest) error {
	if p == nil || !p.CanManageApplications() {
		return ErrPermissionDenied
	}
	return s.repo.UpdateAppConfigurations(ctx, customerID(p), req)
}

func (s *Service) GetVersionConfigurations(ctx context.Context, p *platformauth.Principal, versionID int) ([]domain.ApplicationVersionConfigurationLink, error) {
	if p == nil || !p.CanManageApplications() {
		return nil, ErrPermissionDenied
	}
	return s.repo.GetVersionConfigurations(ctx, customerID(p), versionID)
}

func (s *Service) UpdateVersionConfigurations(ctx context.Context, p *platformauth.Principal, req domain.LinkConfigurationsToAppVersionRequest) error {
	if p == nil || !p.CanManageApplications() {
		return ErrPermissionDenied
	}
	return s.repo.UpdateVersionConfigurations(ctx, customerID(p), req)
}

func (s *Service) AdminSearch(ctx context.Context, p *platformauth.Principal, value string) ([]domain.Application, error) {
	if p == nil || !p.SuperAdmin {
		return nil, ErrPermissionDenied
	}
	return s.repo.AdminSearch(ctx, value)
}

func (s *Service) TurnIntoCommon(ctx context.Context, p *platformauth.Principal, id int) error {
	if p == nil || !p.SuperAdmin {
		return ErrPermissionDenied
	}
	return s.repo.TurnIntoCommon(ctx, id)
}
