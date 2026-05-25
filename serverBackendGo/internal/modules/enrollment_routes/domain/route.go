package domain

import "time"

// EnrollmentRouteView is the admin API DTO (no profile fields).
type EnrollmentRouteView struct {
	ID                            int        `json:"id"`
	Name                          string     `json:"name"`
	Description                   string     `json:"description,omitempty"`
	QRCodeKey                     string     `json:"qrcodekey,omitempty"`
	TargetNodeID                  int        `json:"targetNodeId"`
	TargetNodeName                string     `json:"targetNodeName,omitempty"`
	TargetNodePath                string     `json:"targetNodePath,omitempty"`
	TargetPlacementKind           string     `json:"targetPlacementKind,omitempty"`
	ContainerPlacementAcknowledged bool       `json:"containerPlacementAcknowledged"`
	DeviceIdentityMode            string     `json:"deviceIdentityMode"`
	BootstrapIntent               string     `json:"bootstrapIntent"`
	BootstrapApplicationID        int        `json:"bootstrapApplicationId"`
	BootstrapApplicationName      string     `json:"bootstrapApplicationName,omitempty"`
	BootstrapVersionID            *int       `json:"bootstrapVersionId,omitempty"`
	ResolvedMainAppVersionID      *int       `json:"resolvedMainAppVersionId,omitempty"`
	ResolvedVersionLabel          string     `json:"resolvedVersionLabel,omitempty"`
	ResolvedPackage               string     `json:"resolvedPackage,omitempty"`
	Status                        string     `json:"status"`
	Type                          int        `json:"type,omitempty"`
	ContainerPlacementAckAt       *time.Time `json:"-"`
	// Provisioning settings (used in QR code generation)
	WifiSSID         string `json:"wifiSsid,omitempty"`
	WifiPassword     string `json:"wifiPassword,omitempty"`
	WifiSecurityType string `json:"wifiSecurityType,omitempty"`
	QRParameters     string `json:"qrParameters,omitempty"`
	AdminExtras      string `json:"adminExtras,omitempty"`
	MobileEnrollment bool   `json:"mobileEnrollment"`
	EncryptDevice    bool   `json:"encryptDevice"`
}

// CreateRequest is POST /enrollment-routes body.
type CreateRequest struct {
	Name                          string  `json:"name"`
	Description                   *string `json:"description,omitempty"`
	TargetNodeID                  int     `json:"targetNodeId"`
	DefaultTreeNodeID             int     `json:"defaultTreeNodeId"` // alias accepted from legacy clients
	DeviceIdentityMode            *string `json:"deviceIdentityMode,omitempty"`
	DefaultDeviceIDMode           *string `json:"defaultDeviceIdMode,omitempty"`
	BootstrapIntent               string  `json:"bootstrapIntent"`
	BootstrapApplicationID        int     `json:"bootstrapApplicationId"`
	BootstrapVersionID            *int    `json:"bootstrapVersionId,omitempty"`
	AcknowledgeContainerPlacement bool    `json:"acknowledgeContainerPlacement"`
	Type                          *int    `json:"type,omitempty"`
	// Provisioning settings
	WifiSSID         *string `json:"wifiSsid,omitempty"`
	WifiPassword     *string `json:"wifiPassword,omitempty"`
	WifiSecurityType *string `json:"wifiSecurityType,omitempty"`
	QRParameters     *string `json:"qrParameters,omitempty"`
	AdminExtras      *string `json:"adminExtras,omitempty"`
	MobileEnrollment *bool   `json:"mobileEnrollment,omitempty"`
	EncryptDevice    *bool   `json:"encryptDevice,omitempty"`
}

// UpdateRequest is PUT /enrollment-routes/:id body.
type UpdateRequest struct {
	Name                          *string `json:"name,omitempty"`
	Description                   *string `json:"description,omitempty"`
	TargetNodeID                  *int    `json:"targetNodeId,omitempty"`
	DefaultTreeNodeID             *int    `json:"defaultTreeNodeId,omitempty"`
	DeviceIdentityMode            *string `json:"deviceIdentityMode,omitempty"`
	DefaultDeviceIDMode           *string `json:"defaultDeviceIdMode,omitempty"`
	BootstrapIntent               *string `json:"bootstrapIntent,omitempty"`
	BootstrapApplicationID        *int    `json:"bootstrapApplicationId,omitempty"`
	BootstrapVersionID            *int    `json:"bootstrapVersionId,omitempty"`
	AcknowledgeContainerPlacement *bool   `json:"acknowledgeContainerPlacement,omitempty"`
	// Provisioning settings
	WifiSSID         *string `json:"wifiSsid,omitempty"`
	WifiPassword     *string `json:"wifiPassword,omitempty"`
	WifiSecurityType *string `json:"wifiSecurityType,omitempty"`
	QRParameters     *string `json:"qrParameters,omitempty"`
	AdminExtras      *string `json:"adminExtras,omitempty"`
	MobileEnrollment *bool   `json:"mobileEnrollment,omitempty"`
	EncryptDevice    *bool   `json:"encryptDevice,omitempty"`
}

// TreeNodeOption is GET /options/tree-nodes item.
type TreeNodeOption struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	ParentID       *int   `json:"parentId"`
	PlacementKind  string `json:"placementKind"`
	DeviceCount    int    `json:"deviceCount"`
	HeavilyLoaded  bool   `json:"heavilyLoaded"`
}

// BootstrapAppVersionOption is a version row in bootstrap apps list.
type BootstrapAppVersionOption struct {
	VersionID     int    `json:"versionId"`
	Version       string `json:"version"`
	VersionCode   int    `json:"versionCode"`
	IsRecommended bool   `json:"isRecommended"`
	IsLatest      bool   `json:"isLatest"`
}

// BootstrapAppOption is GET /options/bootstrap-apps item.
type BootstrapAppOption struct {
	ApplicationID int                         `json:"applicationId"`
	Name          string                      `json:"name"`
	Package       string                      `json:"package"`
	Versions      []BootstrapAppVersionOption `json:"versions"`
}

// QRMeta is GET /enrollment-routes/:id/qr response.
type QRMeta struct {
	QRCodeKey                string                 `json:"qrcodekey"`
	DefaultDeviceIDMode      string                 `json:"defaultDeviceIdMode"`
	ResolvedMainAppVersionID *int                   `json:"resolvedMainAppVersionId,omitempty"`
	MainAppPackage           string                 `json:"mainAppPackage,omitempty"`
	MainAppVersion           string                 `json:"mainAppVersion,omitempty"`
	MainAppVersionCode       int                    `json:"mainAppVersionCode,omitempty"`
	TargetNodeID             int                    `json:"targetNodeId"`
	Contract                 map[string]interface{} `json:"contract,omitempty"`
}

// ResolvedBootstrap is the outcome of intent resolution.
type ResolvedBootstrap struct {
	ApplicationID int
	VersionID     int
	Package       string
	VersionLabel  string
	VersionCode   int
}
