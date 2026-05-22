package domain

// LookupItem mirrors com.hmdm.rest.json.LookupItem.
type LookupItem struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
}

// Application mirrors React Application.
type Application struct {
	ID                 *int    `json:"id,omitempty"`
	Name               *string `json:"name,omitempty"`
	Pkg                *string `json:"pkg,omitempty"`
	Version            *string `json:"version,omitempty"`
	VersionCode        *int    `json:"versionCode,omitempty"`
	URL                *string `json:"url,omitempty"`
	URLArmeabi         *string `json:"urlArmeabi,omitempty"`
	URLArm64           *string `json:"urlArm64,omitempty"`
	Type               *string `json:"type,omitempty"`
	ShowIcon           *bool   `json:"showIcon,omitempty"`
	System             *bool   `json:"system,omitempty"`
	Intent             *string `json:"intent,omitempty"`
	CustomerID         *int    `json:"customerId,omitempty"`
	CustomerName       *string `json:"customerName,omitempty"`
	Common             *bool   `json:"common,omitempty"`
	CommonApplication  *bool   `json:"commonApplication,omitempty"`
	LatestVersion      *int    `json:"latestVersion,omitempty"`
	LatestVersionText  *string `json:"latestVersionText,omitempty"`
	UsedVersionID      *int    `json:"usedVersionId,omitempty"`
	FilePath           *string `json:"filePath,omitempty"`
	RunAfterInstall    *bool   `json:"runAfterInstall,omitempty"`
	RunAtBoot          *bool   `json:"runAtBoot,omitempty"`
	SkipVersion        *bool   `json:"skipVersion,omitempty"`
	IconText           *string `json:"iconText,omitempty"`
	IconID             *int    `json:"iconId,omitempty"`
	UseKiosk           *bool   `json:"useKiosk,omitempty"`
	Split              *bool   `json:"split,omitempty"`
	Arch               *string `json:"arch,omitempty"`
	AutoUpdate         *bool   `json:"autoUpdate,omitempty"`
}

// ApplicationVersion mirrors React ApplicationVersion.
type ApplicationVersion struct {
	ID            *int    `json:"id,omitempty"`
	ApplicationID *int    `json:"applicationId,omitempty"`
	Version       *string `json:"version,omitempty"`
	VersionCode   *int    `json:"versionCode,omitempty"`
	URL           *string `json:"url,omitempty"`
	URLArmeabi    *string `json:"urlArmeabi,omitempty"`
	URLArm64      *string `json:"urlArm64,omitempty"`
	FilePath      *string `json:"filePath,omitempty"`
	Split         *bool   `json:"split,omitempty"`
	Arch          *string `json:"arch,omitempty"`
	Action        *int    `json:"action,omitempty"`
	ShowIcon      *bool   `json:"showIcon,omitempty"`
	ScreenOrder   *int    `json:"screenOrder,omitempty"`
	KeyCode       *int    `json:"keyCode,omitempty"`
	Bottom        *bool   `json:"bottom,omitempty"`
	AutoUpdate    *bool   `json:"autoUpdate,omitempty"`
	ApkHash       *string `json:"apkHash,omitempty"`
}

// ValidatePkgRequest is PUT /validatePkg body.
type ValidatePkgRequest struct {
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Pkg  string  `json:"pkg"`
}

// ApplicationConfigurationLink is a configuration linked to an application.
type ApplicationConfigurationLink struct {
	ID                *int    `json:"id,omitempty"`
	ConfigurationID   int     `json:"configurationId"`
	ApplicationID     *int    `json:"applicationId,omitempty"`
	Name              *string `json:"name,omitempty"`
	Action            int     `json:"action"`
	Selected          *bool   `json:"selected,omitempty"`
	Notify            *bool   `json:"notify,omitempty"`
}

// ApplicationVersionConfigurationLink is version-level link row.
type ApplicationVersionConfigurationLink struct {
	ConfigurationID int     `json:"configurationId"`
	Name            *string `json:"name,omitempty"`
	Action          int     `json:"action"`
	Selected        *bool   `json:"selected,omitempty"`
	Notify          *bool   `json:"notify,omitempty"`
}

// LinkConfigurationsToAppRequest is POST /configurations body.
type LinkConfigurationsToAppRequest struct {
	ApplicationID   int                            `json:"applicationId"`
	Configurations  []ApplicationConfigurationLink   `json:"configurations"`
}

// LinkConfigurationsToAppVersionRequest is POST /version/configurations body.
type LinkConfigurationsToAppVersionRequest struct {
	ApplicationVersionID int                                  `json:"applicationVersionId"`
	Configurations       []ApplicationVersionConfigurationLink `json:"configurations"`
}
