package application

import (
	"context"
	"database/sql"

	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// HubService aggregates profile hub read models (019).
type HubService struct {
	repo      port.ProfileRepository
	hub       *postgres.HubRepository
	staleDays int
}

func NewHubService(db *sql.DB, repo port.ProfileRepository, staleDays int) *HubService {
	if staleDays <= 0 {
		staleDays = 30
	}
	return &HubService{repo: repo, hub: postgres.NewHubRepository(db), staleDays: staleDays}
}

func (s *HubService) List(ctx context.Context, p *platformauth.Principal) ([]domain.ProfileListItem, error) {
	if err := requireConfigPerm(p); err != nil {
		return nil, err
	}
	cid := customerID(p)
	items, err := s.repo.List(ctx, cid)
	if err != nil {
		return nil, err
	}
	metrics, err := s.hub.ListRowMetrics(ctx, cid, s.staleDays)
	if err != nil {
		return nil, err
	}
	for i := range items {
		m := metrics[items[i].ID]
		hm := domain.HubMetrics{
			HasPublished:        m.HasPublished,
			Enabled:             items[i].Enabled,
			AssignmentCount:     m.AssignmentCount,
			RolloutFailureCount: m.RolloutFailureCount,
			HasUnpublishedDraft: m.HasUnpublishedDraft,
			StalePublish:        m.StalePublish,
		}
		h, reasons, badges := ComputeHealth(hm)
		items[i].Health = h
		items[i].HealthReasons = reasons
		items[i].Badges = badges
		items[i].AssignmentCount = m.AssignmentCount
		items[i].RolloutFailureCount = m.RolloutFailureCount
	}
	if items == nil {
		items = []domain.ProfileListItem{}
	}
	return items, nil
}

func (s *HubService) Summary(ctx context.Context, p *platformauth.Principal, profileID int) (*domain.ProfileSummary, error) {
	if err := requireConfigPerm(p); err != nil {
		return nil, err
	}
	cid := customerID(p)
	meta, err := s.repo.GetMeta(ctx, cid, profileID)
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, ErrProfileNotFound
	}
	metrics, err := s.hub.ListRowMetrics(ctx, cid, s.staleDays)
	if err != nil {
		return nil, err
	}
	m, ok := metrics[profileID]
	if !ok {
		m = postgres.ListRowMetrics{ProfileID: profileID}
	}
	hm := domain.HubMetrics{
		HasPublished:        m.HasPublished,
		Enabled:             meta.Enabled,
		AssignmentCount:     m.AssignmentCount,
		RolloutFailureCount: m.RolloutFailureCount,
		HasUnpublishedDraft: m.HasUnpublishedDraft,
		StalePublish:        m.StalePublish,
	}
	health, reasons, _ := ComputeHealth(hm)
	rollout, err := s.hub.RolloutSnapshot(ctx, cid, profileID)
	if err != nil {
		return nil, err
	}
	folders, err := s.hub.AssignedFolderNames(ctx, cid, profileID)
	if err != nil {
		return nil, err
	}
	pinned := domain.ProfilePinned{}
	if meta.PublishedVersionID != nil && *meta.PublishedVersionID > 0 {
		pinned, _ = s.hub.PinnedFromPublished(ctx, *meta.PublishedVersionID)
	}
	var publishedCtx *domain.PublishedContext
	if meta.PublishedVersionID != nil && *meta.PublishedVersionID > 0 {
		vnum := 0
		if meta.PublishedVersion != nil {
			vnum = *meta.PublishedVersion
		}
		publishedCtx = &domain.PublishedContext{
			VersionID:      *meta.PublishedVersionID,
			VersionNumber:  vnum,
			Status:         "published",
			PinnedSettings: pinned,
		}
	}
	summary := &domain.ProfileSummary{
		ID:                     profileID,
		Name:                   meta.Name,
		Description:            meta.Description,
		Enabled:                meta.Enabled,
		Health:                 health,
		HealthReasons:          reasons,
		Lifecycle:              LifecycleLabel(meta.Enabled, m.HasPublished),
		PublishedVersionID:     meta.PublishedVersionID,
		PublishedVersionNumber: meta.PublishedVersion,
		DraftVersionID:         meta.DraftVersionID,
		HasUnpublishedDraft:    m.HasUnpublishedDraft,
		CanPublish:             m.HasUnpublishedDraft,
		AssignmentCount:        m.AssignmentCount,
		AssignedFolders:        folders,
		Rollout:                rollout,
		PinnedSettings:         pinned,
		PublishedContext:       publishedCtx,
	}
	return summary, nil
}

func (s *HubService) Activity(ctx context.Context, p *platformauth.Principal, profileID int, limit int) (*domain.ProfileActivityPage, error) {
	if err := requireConfigPerm(p); err != nil {
		return nil, err
	}
	meta, err := s.repo.GetMeta(ctx, customerID(p), profileID)
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, ErrProfileNotFound
	}
	items, err := s.hub.ListActivity(ctx, profileID, limit)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.ProfileActivityEvent{}
	}
	return &domain.ProfileActivityPage{Items: items}, nil
}
