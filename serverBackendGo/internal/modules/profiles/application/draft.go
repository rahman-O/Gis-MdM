package application

import (
	"context"
	"errors"

	cfgdomain "github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// DraftService handles profile draft read/write (US3).
type DraftService struct {
	repo port.ProfileRepository
}

func NewDraftService(repo port.ProfileRepository) *DraftService {
	return &DraftService{repo: repo}
}

var (
	ErrPermissionDenied = errors.New("error.permission.denied")
	ErrProfileNotFound  = errors.New("error.notfound.profile")
	ErrVersionNotFound  = errors.New("error.notfound.profile.version")
	ErrNotDraftVersion  = errors.New("error.profile.version.notdraft")
	ErrDuplicateProfile = errors.New("error.duplicate.profile")
)

func customerID(p *platformauth.Principal) int {
	if p == nil {
		return 0
	}
	return p.CustomerID
}

func (s *DraftService) requireConfigPerm(p *platformauth.Principal) error {
	if p == nil || !p.CanManageConfigurations() {
		return ErrPermissionDenied
	}
	return nil
}

func (s *DraftService) List(ctx context.Context, p *platformauth.Principal) ([]domain.ProfileListItem, error) {
	if err := s.requireConfigPerm(p); err != nil {
		return nil, err
	}
	items, err := s.repo.List(ctx, customerID(p))
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.ProfileListItem{}
	}
	return items, nil
}

func (s *DraftService) GetMeta(ctx context.Context, p *platformauth.Principal, profileID int) (*domain.ProfileMeta, error) {
	if err := s.requireConfigPerm(p); err != nil {
		return nil, err
	}
	meta, err := s.repo.GetMeta(ctx, customerID(p), profileID)
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, ErrProfileNotFound
	}
	if meta.DraftVersionID == nil {
		draftID, err := s.repo.EnsureDraft(ctx, customerID(p), profileID)
		if err != nil && !errors.Is(err, postgres.ErrNoPublishedToFork) {
			return nil, err
		}
		if err == nil {
			meta.DraftVersionID = &draftID
		}
	}
	return meta, nil
}

type VersionResponse struct {
	domain.VersionMeta
	Payload cfgdomain.Configuration `json:",inline"`
}

func (s *DraftService) GetVersion(ctx context.Context, p *platformauth.Principal, profileID, versionID int) (*VersionResponse, error) {
	if err := s.requireConfigPerm(p); err != nil {
		return nil, err
	}
	payload, meta, err := s.repo.GetVersion(ctx, customerID(p), profileID, versionID)
	if err != nil {
		return nil, err
	}
	if payload == nil || meta == nil {
		return nil, ErrVersionNotFound
	}
	return &VersionResponse{VersionMeta: *meta, Payload: *payload}, nil
}

func (s *DraftService) Create(ctx context.Context, p *platformauth.Principal, req domain.CreateRequest) (*domain.ProfileMeta, error) {
	if err := s.requireConfigPerm(p); err != nil {
		return nil, err
	}
	profileID, draftID, err := s.repo.Create(ctx, customerID(p), req)
	if err != nil {
		if errors.Is(err, postgres.ErrDuplicateProfileName) {
			return nil, ErrDuplicateProfile
		}
		return nil, err
	}
	return &domain.ProfileMeta{
		ID: profileID, Name: req.Name,
		Description: derefStr(req.Description),
		DraftVersionID: &draftID,
	}, nil
}

func (s *DraftService) SaveDraft(ctx context.Context, p *platformauth.Principal, profileID, versionID int, payload cfgdomain.Configuration) error {
	if err := s.requireConfigPerm(p); err != nil {
		return err
	}
	err := s.repo.SaveDraft(ctx, customerID(p), profileID, versionID, payload)
	if errors.Is(err, postgres.ErrVersionNotFound) {
		return ErrVersionNotFound
	}
	if errors.Is(err, postgres.ErrNotDraftVersion) {
		return ErrNotDraftVersion
	}
	return err
}

func (s *DraftService) EnsureDraft(ctx context.Context, p *platformauth.Principal, profileID int) (int, error) {
	if err := s.requireConfigPerm(p); err != nil {
		return 0, err
	}
	id, err := s.repo.EnsureDraft(ctx, customerID(p), profileID)
	if errors.Is(err, postgres.ErrProfileNotFound) {
		return 0, ErrProfileNotFound
	}
	return id, err
}

func (s *DraftService) ListVersions(ctx context.Context, p *platformauth.Principal, profileID int) ([]domain.VersionListItem, error) {
	if err := s.requireConfigPerm(p); err != nil {
		return nil, err
	}
	return s.repo.ListVersions(ctx, customerID(p), profileID)
}

func (s *DraftService) ForkDraftFromVersion(ctx context.Context, p *platformauth.Principal, profileID, sourceVersionID int) (*domain.ProfileMeta, error) {
	if err := s.requireConfigPerm(p); err != nil {
		return nil, err
	}
	_, err := s.repo.ForkDraftFromPublished(ctx, customerID(p), profileID, sourceVersionID)
	if errors.Is(err, postgres.ErrVersionNotFound) {
		return nil, ErrVersionNotFound
	}
	if errors.Is(err, postgres.ErrNotDraftVersion) {
		return nil, ErrVersionNotPublished
	}
	if err != nil {
		return nil, err
	}
	return s.GetMeta(ctx, p, profileID)
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
