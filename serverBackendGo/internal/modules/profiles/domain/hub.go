package domain

import "time"

// Profile health values for list and workspace cockpit.
const (
	HealthHealthy   = "healthy"
	HealthWarning   = "warning"
	HealthError     = "error"
	HealthDraftOnly = "draft_only"
)

// ProfileListItemHub fields extend list rows (019).
type ProfileListHubFields struct {
	Health              string   `json:"health"`
	HealthReasons       []string `json:"healthReasons,omitempty"`
	Badges              []string `json:"badges,omitempty"`
	AssignmentCount     int      `json:"assignmentCount"`
	RolloutFailureCount int      `json:"rolloutFailureCount"`
}

// PublishedContext is the published policy strip for Assignments (020).
type PublishedContext struct {
	VersionID       int           `json:"versionId"`
	VersionNumber   int           `json:"versionNumber"`
	Status          string        `json:"status"`
	PinnedSettings  ProfilePinned `json:"pinnedSettings"`
}

// ProfileSummary is GET /profiles/:id/summary.
type ProfileSummary struct {
	ID                    int                `json:"id"`
	Name                  string             `json:"name"`
	Description           string             `json:"description"`
	Enabled               bool               `json:"enabled"`
	Health                string             `json:"health"`
	HealthReasons         []string           `json:"healthReasons,omitempty"`
	Lifecycle             string             `json:"lifecycle"`
	PublishedVersionID    *int               `json:"publishedVersionId,omitempty"`
	PublishedVersionNumber *int              `json:"publishedVersionNumber,omitempty"`
	DraftVersionID        *int               `json:"draftVersionId,omitempty"`
	HasUnpublishedDraft   bool               `json:"hasUnpublishedDraft"`
	CanPublish            bool               `json:"canPublish"`
	AssignmentCount       int                `json:"assignmentCount"`
	AssignedFolders       []string           `json:"assignedFolders"`
	Rollout               ProfileRolloutSnap `json:"rollout"`
	PinnedSettings        ProfilePinned      `json:"pinnedSettings"`
	PublishedContext      *PublishedContext  `json:"publishedContext"`
}

// ProfileRolloutSnap is rollout counts for overview cards.
type ProfileRolloutSnap struct {
	Pending   int `json:"pending"`
	Installed int `json:"installed"`
	Partial   int `json:"partial"`
	Failed    int `json:"failed"`
	Total     int `json:"total"`
}

// ProfilePinned is read-only highlights for overview.
type ProfilePinned struct {
	KioskMode       bool       `json:"kioskMode"`
	MainAppName     string     `json:"mainAppName,omitempty"`
	AppCount        int        `json:"appCount"`
	LastPublishedAt *time.Time `json:"lastPublishedAt,omitempty"`
}

// ProfileActivityEvent is one timeline row.
type ProfileActivityEvent struct {
	ID            int64          `json:"id"`
	EventType     string         `json:"eventType"`
	SummaryKey    string         `json:"summaryKey"`
	SummaryParams map[string]any `json:"summaryParams,omitempty"`
	OccurredAt    time.Time      `json:"occurredAt"`
	ActorUserID   *int           `json:"actorUserId,omitempty"`
}

// ProfileActivityPage is GET /profiles/:id/activity.
type ProfileActivityPage struct {
	Items []ProfileActivityEvent `json:"items"`
}

// HubMetrics is internal input for health computation.
type HubMetrics struct {
	HasPublished        bool
	Enabled             bool
	AssignmentCount     int
	RolloutFailureCount int
	HasUnpublishedDraft bool
	StalePublish        bool
}
