package domain

// Rollout status values stored on devices.profile_rollout_status.
const (
	RolloutPending   = "pending"
	RolloutInstalling = "installing"
	RolloutInstalled = "installed"
	RolloutPartial   = "partial"
	RolloutFailed    = "failed"
)

// VersionListItem is one row in GET /profiles/:id/versions.
type VersionListItem struct {
	VersionID     int     `json:"versionId"`
	VersionNumber int     `json:"versionNumber"`
	Status        string  `json:"status"`
	PublishedAt   *string `json:"publishedAt,omitempty"`
	CreatedAt     string  `json:"createdAt"`
}

// DeviceRolloutRow is one device in rollout status grid.
type DeviceRolloutRow struct {
	DeviceID             int     `json:"deviceId"`
	DeviceName           string  `json:"deviceName"`
	TreeNodeID           *int    `json:"treeNodeId,omitempty"`
	TreeNodeName         string  `json:"treeNodeName,omitempty"`
	TargetVersionID      *int    `json:"targetVersionId,omitempty"`
	TargetVersionNumber  *int    `json:"targetVersionNumber,omitempty"`
	AppliedVersionID     *int    `json:"appliedVersionId,omitempty"`
	AppliedVersionNumber *int    `json:"appliedVersionNumber,omitempty"`
	Status               string  `json:"status"`
	Reason               string  `json:"reason,omitempty"`
	LastUpdate           *int64  `json:"lastUpdate,omitempty"`
	ResolutionSource     string  `json:"resolutionSource,omitempty"`
}

// RolloutDevicesQuery filters GET rollout/devices.
type RolloutDevicesQuery struct {
	TreeNodeID *int
	Status     string
	Page       int
	PageSize   int
}

// RolloutDevicesPage is paginated rollout list.
type RolloutDevicesPage struct {
	Items      []DeviceRolloutRow `json:"items"`
	TotalCount int                `json:"totalCount"`
}

// EffectiveProfileResolution is the resolved policy target for a device.
type EffectiveProfileResolution struct {
	ProfileID        int
	ProfileVersionID int
	RouteID          int64
	Source           string // tree | route | none
	Enabled          bool
}

// EnableProfileResult is disable/enable API data.
type EnableProfileResult struct {
	Enabled              bool `json:"enabled"`
	DevicesMarkedPending int  `json:"devicesMarkedPending"`
}
