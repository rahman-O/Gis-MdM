package application

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

const eventProfileVersionDeleted = "ProfileVersionDeleted"

var (
	ErrVersionDeleteActivePublished = errors.New("error.profile.version.delete.activePublished")
	ErrVersionDeleteAssigned        = errors.New("error.profile.version.delete.assigned")
	ErrVersionDeleteDevicesTarget   = errors.New("error.profile.version.delete.devicesTarget")
)

// VersionDeleteService removes draft or unused historical versions (020).
type VersionDeleteService struct {
	repo port.ProfileRepository
	db   *sql.DB
}

func NewVersionDeleteService(repo port.ProfileRepository, db *sql.DB) *VersionDeleteService {
	return &VersionDeleteService{repo: repo, db: db}
}

func (s *VersionDeleteService) Delete(ctx context.Context, p *platformauth.Principal, profileID, versionID int) (*domain.VersionDeleteResult, error) {
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
	active, assigned, devicesTarget, err := s.repo.VersionDeleteEligibility(ctx, customerID(p), profileID, versionID)
	if err != nil {
		if errors.Is(err, postgres.ErrVersionNotFound) {
			return nil, ErrVersionNotFound
		}
		if errors.Is(err, postgres.ErrProfileNotFound) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	if active {
		return nil, ErrVersionDeleteActivePublished
	}
	if assigned {
		return nil, ErrVersionDeleteAssigned
	}
	if devicesTarget {
		return nil, ErrVersionDeleteDevicesTarget
	}
	if err := s.repo.DeleteVersion(ctx, customerID(p), profileID, versionID); err != nil {
		if errors.Is(err, postgres.ErrVersionNotFound) {
			return nil, ErrVersionNotFound
		}
		if errors.Is(err, postgres.ErrProfileNotFound) {
			return nil, ErrProfileNotFound
		}
		if errors.Is(err, postgres.ErrVersionDeleteActivePublished) {
			return nil, ErrVersionDeleteActivePublished
		}
		if errors.Is(err, postgres.ErrVersionDeleteAssigned) {
			return nil, ErrVersionDeleteAssigned
		}
		if errors.Is(err, postgres.ErrVersionDeleteDevicesTarget) {
			return nil, ErrVersionDeleteDevicesTarget
		}
		return nil, err
	}
	payload, _ := json.Marshal(map[string]int{"profileId": profileID, "versionId": versionID})
	_ = s.repo.InsertDomainEvent(ctx, eventProfileVersionDeleted, fmt.Sprintf("profile:%d", profileID), payload)
	return &domain.VersionDeleteResult{ProfileID: profileID, VersionID: versionID}, nil
}
