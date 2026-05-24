package domain

import cfgdomain "github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"

// VersionPayload is the full profile version editor body (same JSON shape as legacy Configuration).
type VersionPayload = cfgdomain.Configuration

// ProfileListItem is returned by GET /profiles.
type ProfileListItem struct {
	ID                   int      `json:"id"`
	Name                 string   `json:"name"`
	Description          string   `json:"description"`
	Enabled              bool     `json:"enabled"`
	Health               string   `json:"health,omitempty"`
	HealthReasons        []string `json:"healthReasons,omitempty"`
	Badges               []string `json:"badges,omitempty"`
	PublishedVersion     *int     `json:"publishedVersion,omitempty"`
	DraftVersionID       *int     `json:"draftVersionId,omitempty"`
	AssignmentCount      int      `json:"assignmentCount"`
	RolloutFailureCount  int      `json:"rolloutFailureCount"`
	DeviceCount          int      `json:"deviceCount"`
	EnrollmentRouteCount int      `json:"enrollmentRouteCount"`
}

// ProfileMeta is returned by GET /profiles/:id.
type ProfileMeta struct {
	ID                   int     `json:"id"`
	Name                 string  `json:"name"`
	Description          string  `json:"description"`
	Enabled              bool    `json:"enabled"`
	DraftVersionID       *int    `json:"draftVersionId,omitempty"`
	PublishedVersionID   *int    `json:"publishedVersionId,omitempty"`
	PublishedVersion     *int    `json:"publishedVersion,omitempty"`
	DeviceCount          int     `json:"deviceCount"`
	EnrollmentRouteCount int     `json:"enrollmentRouteCount"`
}

// VersionMeta is editor metadata alongside the payload.
type VersionMeta struct {
	ProfileID     int    `json:"profileId"`
	VersionID     int    `json:"versionId"`
	VersionNumber int    `json:"versionNumber"`
	Status        string `json:"status"`
}

// CreateRequest is POST /profiles body.
type CreateRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}
