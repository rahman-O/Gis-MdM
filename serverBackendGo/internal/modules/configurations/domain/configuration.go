package domain

// LookupItem mirrors com.hmdm.rest.json.LookupItem.
type LookupItem struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
}

// ConfigurationApplication is a row on the configuration Applications tab.
type ConfigurationApplication struct {
	ID              int     `json:"id"`
	Name            *string `json:"name,omitempty"`
	Pkg             *string `json:"pkg,omitempty"`
	Type            *string `json:"type,omitempty"`
	Action          *int    `json:"action,omitempty"`
	Version         *string `json:"version,omitempty"`
	URL             *string `json:"url,omitempty"`
	VersionCode     *int    `json:"versionCode,omitempty"`
	UsedVersionID   *int    `json:"usedVersionId,omitempty"`
	LatestVersion   *int    `json:"latestVersion,omitempty"`
	ShowIcon        *bool   `json:"showIcon,omitempty"`
	ScreenOrder     *int    `json:"screenOrder,omitempty"`
	KeyCode         *int    `json:"keyCode,omitempty"`
	Bottom            *bool `json:"bottom,omitempty"`
	Remove            *bool `json:"remove,omitempty"`
	LongTap           *bool `json:"longTap,omitempty"`
	SkipVersionCheck  *bool `json:"skipVersionCheck,omitempty"`
}

// ConfigurationFile is a file entry on a configuration.
type ConfigurationFile struct {
	ID          *int    `json:"id,omitempty"`
	Path        *string `json:"path,omitempty"`
	ExternalURL *string `json:"externalUrl,omitempty"`
	URL         *string `json:"url,omitempty"`
	Remove      *bool   `json:"remove,omitempty"`
}

// ConfigurationApplicationSetting is per-app override on a configuration.
type ConfigurationApplicationSetting struct {
	ID              *int    `json:"id,omitempty"`
	ApplicationID   *int    `json:"applicationId,omitempty"`
	ApplicationName *string `json:"applicationName,omitempty"`
	Name            *string `json:"name,omitempty"`
	Type            *string `json:"type,omitempty"`
	Value           *string `json:"value,omitempty"`
}

// Configuration mirrors React Configuration (editor + list).
type Configuration struct {
	ID                       *int                              `json:"id,omitempty"`
	Name                     *string                           `json:"name,omitempty"`
	Description              *string                           `json:"description,omitempty"`
	Type                     *int                              `json:"type,omitempty"`
	DeviceCount              *int                              `json:"deviceCount,omitempty"`
	Password                 *string                           `json:"password,omitempty"`
	BackgroundColor          *string                           `json:"backgroundColor,omitempty"`
	TextColor                *string                           `json:"textColor,omitempty"`
	BackgroundImageURL       *string                           `json:"backgroundImageUrl,omitempty"`
	QRCodeKey                *string                           `json:"qrCodeKey,omitempty"`
	BaseURL                  *string                           `json:"baseUrl,omitempty"`
	MainAppID                *int                              `json:"mainAppId,omitempty"`
	ContentAppID             *int                              `json:"contentAppId,omitempty"`
	DefaultFilePath          *string                           `json:"defaultFilePath,omitempty"`
	Permissive               *bool                             `json:"permissive,omitempty"`
	Applications             []ConfigurationApplication        `json:"applications,omitempty"`
	Files                    []ConfigurationFile               `json:"files,omitempty"`
	ApplicationSettings      []ConfigurationApplicationSetting `json:"applicationSettings,omitempty"`
	// Policy holds MDM keys from settingsjson (merged into API JSON on read/write).
	Policy map[string]any `json:"-"`
}

// CopyRequest is PUT /copy body.
type CopyRequest struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// UpgradeApplicationRequest is PUT /application/upgrade body.
type UpgradeApplicationRequest struct {
	ConfigurationID int `json:"configurationId"`
	ApplicationID   int `json:"applicationId"`
}
