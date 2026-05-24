package domain

// Route is an enrollment route (QR + tree placement + published profile version).
type Route struct {
	ID                   int     `json:"id"`
	Name                 string  `json:"name"`
	Description          string  `json:"description,omitempty"`
	QRCodeKey            string  `json:"qrcodekey,omitempty"`
	ProfileID            int     `json:"profileId,omitempty"`
	ProfileVersionID     int     `json:"profileVersionId,omitempty"`
	ProfileVersionNumber *int    `json:"profileVersionNumber,omitempty"`
	DefaultTreeNodeID    int     `json:"defaultTreeNodeId,omitempty"`
	DefaultTreeNodeName  string  `json:"defaultTreeNodeName,omitempty"`
	DefaultDeviceIDMode  string  `json:"defaultDeviceIdMode"`
	MainAppID            *int    `json:"mainAppId,omitempty"`
	Type                 int     `json:"type,omitempty"`
}

// RouteDetail is the editor payload.
type RouteDetail struct {
	Route
}

// CreateRequest is POST /enrollment-routes body.
type CreateRequest struct {
	Name                string  `json:"name"`
	Description         *string `json:"description,omitempty"`
	ProfileVersionID    *int    `json:"profileVersionId,omitempty"`
	DefaultTreeNodeID   int     `json:"defaultTreeNodeId"`
	DefaultDeviceIDMode *string `json:"defaultDeviceIdMode,omitempty"`
	MainAppID           *int    `json:"mainAppId,omitempty"`
	Type                *int    `json:"type,omitempty"`
}

// UpdateRequest is PUT /enrollment-routes/:id body.
type UpdateRequest struct {
	Name                *string `json:"name,omitempty"`
	Description         *string `json:"description,omitempty"`
	ProfileVersionID    *int    `json:"profileVersionId,omitempty"`
	DefaultTreeNodeID   *int    `json:"defaultTreeNodeId,omitempty"`
	DefaultDeviceIDMode *string `json:"defaultDeviceIdMode,omitempty"`
	MainAppID           *int    `json:"mainAppId,omitempty"`
}

// PublishedProfileVersion is a picker option for route binding.
type PublishedProfileVersion struct {
	ProfileVersionID   int    `json:"profileVersionId"`
	ProfileID          int    `json:"profileId"`
	ProfileName        string `json:"profileName"`
	VersionNumber      int    `json:"versionNumber"`
	ProfileEnabled     bool   `json:"profileEnabled"`
	MainAppID          *int   `json:"mainAppId,omitempty"`
}

// QRMeta is GET /enrollment-routes/:id/qr response.
type QRMeta struct {
	QRCodeKey           string `json:"qrcodekey"`
	DefaultDeviceIDMode string `json:"defaultDeviceIdMode"`
	MainAppID           *int   `json:"mainAppId,omitempty"`
}
