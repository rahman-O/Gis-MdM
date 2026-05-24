package application

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"

	profilepostgres "github.com/gis-mdm/server-backend-go/internal/modules/profiles/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

const eventProfilePublished = "ProfilePublished"

// PublishService handles profile publish and impact preview.
type PublishService struct {
	repo     port.ProfileRepository
	rollout  port.RolloutStore
	db       *sql.DB
	draft    *DraftService
	compiler *ArtifactCompiler
}

func NewPublishService(repo port.ProfileRepository, rollout port.RolloutStore, db *sql.DB, draft *DraftService, compiler *ArtifactCompiler) *PublishService {
	return &PublishService{repo: repo, rollout: rollout, db: db, draft: draft, compiler: compiler}
}

var ErrConfirmImpactRequired = errors.New("error.profile.publish.confirm_required")

func (s *PublishService) Impact(ctx context.Context, p *platformauth.Principal, profileID int) (*domain.ImpactSummary, error) {
	if err := s.draft.requireConfigPerm(p); err != nil {
		return nil, err
	}
	if _, err := s.draft.GetMeta(ctx, p, profileID); err != nil {
		return nil, err
	}
	cid := customerID(p)
	devices, routes, err := s.repo.CountImpact(ctx, cid, profileID)
	if err != nil {
		return nil, err
	}
	assignments, err := s.listAssignmentImpact(ctx, cid, profileID)
	if err != nil {
		return nil, err
	}
	summary := BuildImpactSummary(devices, routes, assignments)
	return &summary, nil
}

func (s *PublishService) listAssignmentImpact(ctx context.Context, customerID, profileID int) ([]domain.PublishImpactAssignment, error) {
	if s.rollout == nil {
		return []domain.PublishImpactAssignment{}, nil
	}
	return s.rollout.ListAssignmentsForPublishImpact(ctx, customerID, profileID)
}

func (s *PublishService) Publish(ctx context.Context, p *platformauth.Principal, profileID, versionID int, req domain.PublishRequest) (*domain.PublishResult, error) {
	if err := s.draft.requireConfigPerm(p); err != nil {
		return nil, err
	}
	cid := customerID(p)
	devices, routes, err := s.repo.CountImpact(ctx, cid, profileID)
	if err != nil {
		return nil, err
	}
	assignments, err := s.listAssignmentImpact(ctx, cid, profileID)
	if err != nil {
		return nil, err
	}
	impact := BuildImpactSummary(devices, routes, assignments)
	if impact.RequiresConfirmDialog && (req.ConfirmImpact == nil || !*req.ConfirmImpact) {
		return nil, ErrConfirmImpactRequired
	}
	ver, err := s.draft.GetVersion(ctx, p, profileID, versionID)
	if err != nil {
		return nil, err
	}
	if ver.Status != "draft" {
		return nil, ErrNotDraftVersion
	}
	settings, err := ver.Payload.BuildSettingsJSON()
	if err != nil {
		return nil, err
	}
	artifact, hash, err := s.compiler.Compile(ctx, CompileInput{
		ProfileID: profileID, ProfileVersionID: versionID,
		VersionNumber: ver.VersionNumber, Payload: ver.Payload, SettingsJSON: settings,
	})
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(artifact)
	if err != nil {
		return nil, err
	}
	publishedBy := 0
	if p != nil && p.ID > 0 {
		publishedBy = int(p.ID)
	}
	assignUpdated, rolloutDevices, err := s.publishWithAssignmentBump(ctx, cid, profileID, versionID, raw, hash, publishedBy)
	if err != nil {
		if errors.Is(err, profilepostgres.ErrNotDraftVersion) {
			return nil, ErrNotDraftVersion
		}
		if errors.Is(err, profilepostgres.ErrVersionNotFound) {
			return nil, ErrVersionNotFound
		}
		return nil, err
	}
	payload, _ := json.Marshal(map[string]any{
		"profileId": profileID, "profileVersionId": versionID, "assignmentsUpdated": assignUpdated,
	})
	_ = s.repo.InsertDomainEvent(ctx, eventProfilePublished, strconv.Itoa(profileID), payload)

	affectedDevices := devices
	if rolloutDevices > affectedDevices {
		affectedDevices = rolloutDevices
	}
	return &domain.PublishResult{
		PublishedVersionID: versionID,
		VersionNumber:      ver.VersionNumber,
		ArtifactHash:       hash,
		AffectedDevices:    affectedDevices,
		AffectedRoutes:     routes,
		AssignmentsUpdated: assignUpdated,
	}, nil
}

func (s *PublishService) publishWithAssignmentBump(ctx context.Context, customerID, profileID, versionID int, artifactJSON []byte, artifactHash string, publishedBy int) (assignmentsUpdated, devicesAffected int, err error) {
	if s.db == nil {
		if err := s.repo.PublishVersion(ctx, customerID, profileID, versionID, artifactJSON, artifactHash, publishedBy); err != nil {
			return 0, 0, err
		}
		if s.rollout == nil {
			return 0, 0, nil
		}
		return s.rollout.BumpAllAssignmentsOnPublish(ctx, customerID, profileID, versionID)
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()
	repo, ok := s.repo.(*profilepostgres.ProfileRepository)
	if !ok {
		return 0, 0, errors.New("publish transaction requires postgres profile repository")
	}
	if err := repo.PublishVersionTx(ctx, tx, customerID, profileID, versionID, artifactJSON, artifactHash, publishedBy); err != nil {
		return 0, 0, err
	}
	if s.rollout != nil {
		assignRepo, ok := s.rollout.(*profilepostgres.AssignmentRepository)
		if !ok {
			return 0, 0, errors.New("publish transaction requires postgres assignment repository")
		}
		assignmentsUpdated, devicesAffected, err = assignRepo.BumpAllAssignmentsOnPublishTx(ctx, tx, customerID, profileID, versionID)
		if err != nil {
			return 0, 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}
	return assignmentsUpdated, devicesAffected, nil
}
