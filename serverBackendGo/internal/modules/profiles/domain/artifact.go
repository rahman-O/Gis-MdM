package domain

import (
	"encoding/json"

	syncdomain "github.com/gis-mdm/server-backend-go/internal/modules/sync/domain"
)

// ProfileArtifact is the frozen policy document produced at publish time.
type ProfileArtifact struct {
	ProfileID        int    `json:"profileId"`
	ProfileVersionID int    `json:"profileVersionId"`
	VersionNumber    int    `json:"versionNumber"`
	Password         string `json:"password,omitempty"`
	BackgroundColor  *string `json:"backgroundColor,omitempty"`
	TextColor        *string `json:"textColor,omitempty"`
	BackgroundImageURL *string `json:"backgroundImageUrl,omitempty"`
	Permissive       bool   `json:"permissive"`
	SettingsJSON     json.RawMessage `json:"settingsJson"`
	Applications     []syncdomain.SyncApplication `json:"applications"`
	Files            []syncdomain.SyncConfigurationFile `json:"files"`
	ConfigApplicationSettings []syncdomain.SyncApplicationSetting `json:"configApplicationSettings,omitempty"`
}

// ImpactSummary is returned by GET /profiles/:id/impact.
type ImpactSummary struct {
	DeviceCount             int                       `json:"deviceCount"`
	EnrollmentRouteCount    int                       `json:"enrollmentRouteCount"`
	RequiresConfirmDialog   bool                      `json:"requiresConfirmDialog"`
	AssignmentsToUpdate     []PublishImpactAssignment `json:"assignmentsToUpdate,omitempty"`
}

// PublishRequest is POST .../publish body.
type PublishRequest struct {
	ConfirmImpact *bool `json:"confirmImpact,omitempty"`
}

// PublishResult is the publish response.
type PublishResult struct {
	PublishedVersionID  int    `json:"publishedVersionId"`
	VersionNumber       int    `json:"versionNumber"`
	ArtifactHash        string `json:"artifactHash"`
	AffectedDevices     int    `json:"affectedDevices"`
	AffectedRoutes      int    `json:"affectedRoutes"`
	AssignmentsUpdated  int    `json:"assignmentsUpdated"`
}
