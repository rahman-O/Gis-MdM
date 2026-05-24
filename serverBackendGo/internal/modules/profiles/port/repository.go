package port

import (
	"context"

	cfgdomain "github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
)

// ProfileRepository persists profiles and versioned policy payloads.
type ProfileRepository interface {
	List(ctx context.Context, customerID int) ([]domain.ProfileListItem, error)
	GetMeta(ctx context.Context, customerID, profileID int) (*domain.ProfileMeta, error)
	GetVersion(ctx context.Context, customerID, profileID, versionID int) (*cfgdomain.Configuration, *domain.VersionMeta, error)
	Create(ctx context.Context, customerID int, req domain.CreateRequest) (profileID, draftVersionID int, err error)
	EnsureDraft(ctx context.Context, customerID, profileID int) (draftVersionID int, err error)
	SaveDraft(ctx context.Context, customerID, profileID, versionID int, payload cfgdomain.Configuration) error
	ListVersionApplications(ctx context.Context, customerID, versionID int) ([]cfgdomain.ConfigurationApplication, error)
	CountImpact(ctx context.Context, customerID, profileID int) (deviceCount, routeCount int, err error)
	PublishVersion(ctx context.Context, customerID, profileID, versionID int, artifactJSON []byte, artifactHash string, publishedBy int) error
	InsertDomainEvent(ctx context.Context, eventType, aggregateID string, payload []byte) error
	ListVersions(ctx context.Context, customerID, profileID int) ([]domain.VersionListItem, error)
	ForkDraftFromPublished(ctx context.Context, customerID, profileID, sourceVersionID int) (int, error)
	DeleteVersion(ctx context.Context, customerID, profileID, versionID int) error
	VersionDeleteEligibility(ctx context.Context, customerID, profileID, versionID int) (activePublished, assigned, devicesTarget bool, err error)
}
