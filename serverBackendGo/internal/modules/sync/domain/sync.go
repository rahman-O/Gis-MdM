package domain

// DeviceCreateOptions mirrors enrollment POST body.
type DeviceCreateOptions struct {
	Customer       string   `json:"customer"`
	Configuration  string   `json:"configuration"`
	Groups         []string `json:"groups"`
}

// DeviceInfo mirrors agent telemetry POST body.
type DeviceInfo struct {
	DeviceID     string          `json:"deviceId"`
	BatteryLevel *int            `json:"batteryLevel"`
	IMEI         *string         `json:"imei"`
	Custom1      *string         `json:"custom1"`
	Custom2      *string         `json:"custom2"`
	Custom3      *string         `json:"custom3"`
	Location     *DeviceLocation `json:"location"`
}

type DeviceLocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Ts  int64   `json:"ts"`
}

// SyncApplicationSetting is persisted per device/package/name.
type SyncApplicationSetting struct {
	PackageID  string `json:"packageId"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
	Value      string `json:"value"`
	Readonly   bool   `json:"readonly"`
	LastUpdate int64  `json:"lastUpdate"`
}

// SyncApplication is an app row in SyncResponse.
type SyncApplication struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Pkg     string  `json:"pkg"`
	Version string  `json:"version"`
	URL     string  `json:"url"`
	Type    string  `json:"type"`
	Icon    *string `json:"icon,omitempty"`
}

// SyncConfigurationFile is a file entry in SyncResponse.
type SyncConfigurationFile struct {
	Path       string `json:"devicePath"`
	URL        string `json:"url"`
	Remove     bool   `json:"remove"`
	External   bool   `json:"external,omitempty"`
}

// SyncResponse is the agent configuration payload.
type SyncResponse struct {
	DeviceID          string                  `json:"deviceId"`
	ConfigurationID   int64                   `json:"configurationId"`
	Password          string                  `json:"password,omitempty"`
	BackgroundColor   *string                 `json:"backgroundColor,omitempty"`
	TextColor         *string                 `json:"textColor,omitempty"`
	Applications      []SyncApplication       `json:"applications"`
	Files             []SyncConfigurationFile `json:"files"`
	ApplicationSettings []SyncApplicationSetting `json:"applicationSettings,omitempty"`
	PushOptions       *string                 `json:"pushOptions,omitempty"`
	RequestUpdates    *string                 `json:"requestUpdates,omitempty"`
	Permissive        *bool                   `json:"permissive,omitempty"`
	AppName           *string                 `json:"appName,omitempty"`
	Vendor            *string                 `json:"vendor,omitempty"`
	NewNumber         *string                 `json:"newNumber,omitempty"`
	Custom1           *string                 `json:"custom1,omitempty"`
	Custom2           *string                 `json:"custom2,omitempty"`
	Custom3           *string                 `json:"custom3,omitempty"`
}

// DeviceRecord is internal persistence shape.
type DeviceRecord struct {
	ID              int64
	CustomerID      int64
	ConfigurationID int64
	Number          string
	OldNumber       *string
	IMEI            *string
	Phone           *string
	LastUpdate      int64
	Custom1         *string
	Custom2         *string
	Custom3         *string
	Info            *string
}
