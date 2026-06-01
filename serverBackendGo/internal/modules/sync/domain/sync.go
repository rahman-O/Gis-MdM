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
	// Extended telemetry fields (sent by Flutter MDM Agent)
	Model          *string `json:"model"`
	Manufacturer   *string `json:"manufacturer"`
	AndroidVersion *string `json:"androidVersion"`
	Serial         *string `json:"serial"`
	Phone          *string `json:"phone"`
	PublicIp       *string `json:"publicIp"`
	MdmMode        *bool   `json:"mdmMode"`
	KioskMode      *bool   `json:"kioskMode"`
	LauncherVersion *string `json:"launcherVersion"`
	DefaultLauncher *string `json:"defaultLauncher"`
	EnrollTime     *int64  `json:"enrollTime"`
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

// SyncApplication is an app row in SyncResponse (Headwind launcher / SyncApplicationInt parity).
type SyncApplication struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Pkg             string  `json:"pkg"`
	Version         string  `json:"version"`
	URL             string  `json:"url"`
	Type            string  `json:"type"`
	Code            *int    `json:"code,omitempty"`
	Icon            *string `json:"icon,omitempty"`
	ShowIcon        *bool   `json:"showIcon,omitempty"`
	ScreenOrder     *int    `json:"screenOrder,omitempty"`
	UseKiosk        *bool   `json:"useKiosk,omitempty"`
	Remove          *bool   `json:"remove,omitempty"`
	System          *bool   `json:"system,omitempty"`
	RunAfterInstall *bool   `json:"runAfterInstall,omitempty"`
	RunAtBoot       *bool   `json:"runAtBoot,omitempty"`
	SkipVersion     *bool   `json:"skipVersion,omitempty"`
	IconText        *string `json:"iconText,omitempty"`
	KeyCode         *int    `json:"keyCode,omitempty"`
	Bottom          *bool   `json:"bottom,omitempty"`
	LongTap         *bool   `json:"longTap,omitempty"`
	Intent          *string `json:"intent,omitempty"`
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
	DeviceID            string                   `json:"deviceId"`
	ConfigurationID     int64                    `json:"configurationId"`
	ProfileID           *int64                   `json:"profileId,omitempty"`
	ProfileVersionID    *int64                   `json:"profileVersionId,omitempty"`
	ProfileRevision     *string                  `json:"profileRevision,omitempty"`
	Password            string                   `json:"password,omitempty"`
	BackgroundColor     *string                  `json:"backgroundColor,omitempty"`
	TextColor           *string                  `json:"textColor,omitempty"`
	BackgroundImageURL  *string                  `json:"backgroundImageUrl,omitempty"`
	Applications        []SyncApplication        `json:"applications"`
	Files               []SyncConfigurationFile  `json:"files"`
	ApplicationSettings []SyncApplicationSetting `json:"applicationSettings,omitempty"`
	PushOptions         *string                  `json:"pushOptions,omitempty"`
	RequestUpdates      *string                  `json:"requestUpdates,omitempty"`
	Permissive          *bool                    `json:"permissive,omitempty"`
	AppName             *string                  `json:"appName,omitempty"`
	Vendor              *string                  `json:"vendor,omitempty"`
	NewNumber           *string                  `json:"newNumber,omitempty"`
	Custom1             *string                  `json:"custom1,omitempty"`
	Custom2             *string                  `json:"custom2,omitempty"`
	Custom3             *string                  `json:"custom3,omitempty"`

	// MDM policy (from configurations.settingsjson — Java SyncResponse parity)
	GPS                 *bool   `json:"gps,omitempty"`
	Bluetooth           *bool   `json:"bluetooth,omitempty"`
	Wifi                *bool   `json:"wifi,omitempty"`
	MobileData          *bool   `json:"mobileData,omitempty"`
	UsbStorage          *bool   `json:"usbStorage,omitempty"`
	LockStatusBar       *bool   `json:"lockStatusBar,omitempty"`
	DisableLocation     *bool   `json:"disableLocation,omitempty"`
	KioskMode           bool    `json:"kioskMode"`
	Restrictions        *string `json:"restrictions,omitempty"`
	KioskHome           *bool   `json:"kioskHome,omitempty"`
	KioskRecents        *bool   `json:"kioskRecents,omitempty"`
	KioskNotifications  *bool   `json:"kioskNotifications,omitempty"`
	KioskSystemInfo     *bool   `json:"kioskSystemInfo,omitempty"`
	KioskKeyguard       *bool   `json:"kioskKeyguard,omitempty"`
	KioskLockButtons    *bool   `json:"kioskLockButtons,omitempty"`
	KioskScreenOn       *bool   `json:"kioskScreenOn,omitempty"`
	KioskExit           *string `json:"kioskExit,omitempty"`
	ShowWifi            *bool   `json:"showWifi,omitempty"`
	LockSafeSettings    *bool   `json:"lockSafeSettings,omitempty"`
	LockVolume          *bool   `json:"lockVolume,omitempty"`
	AutoBrightness      *bool   `json:"autoBrightness,omitempty"`
	ManageTimeout       *bool   `json:"manageTimeout,omitempty"`
	ManageVolume        *bool   `json:"manageVolume,omitempty"`
	RunDefaultLauncher  *bool   `json:"runDefaultLauncher,omitempty"`
	DisableScreenshots  *bool   `json:"disableScreenshots,omitempty"`
	AutostartForeground *bool   `json:"autostartForeground,omitempty"`
	ScheduleAppUpdate   bool    `json:"scheduleAppUpdate,omitempty"`
	AppPermissions      *string `json:"appPermissions,omitempty"`
	PasswordMode        *string `json:"passwordMode,omitempty"`
	TimeZone            *string `json:"timeZone,omitempty"`
	AllowedClasses      *string `json:"allowedClasses,omitempty"`
	NewServerUrl        *string `json:"newServerUrl,omitempty"`
	MainApp             *string `json:"mainApp,omitempty"`
	SystemUpdateType    *int    `json:"systemUpdateType,omitempty"`
	Orientation         *int    `json:"orientation,omitempty"`
	Brightness          *int    `json:"brightness,omitempty"`
	Timeout             *int    `json:"timeout,omitempty"`
	Volume              *int    `json:"volume,omitempty"`
	KeepaliveTime       *int    `json:"keepaliveTime,omitempty"`
	IconSize            *int    `json:"iconSize,omitempty"`
	SystemUpdateFrom    *string `json:"systemUpdateFrom,omitempty"`
	SystemUpdateTo      *string `json:"systemUpdateTo,omitempty"`
	AppUpdateFrom       *string `json:"appUpdateFrom,omitempty"`
	AppUpdateTo         *string `json:"appUpdateTo,omitempty"`
	DownloadUpdates     *string `json:"downloadUpdates,omitempty"`
}

// DeviceRecord is internal persistence shape.
type DeviceRecord struct {
	ID                int64
	CustomerID        int64
	ConfigurationID   int64
	EnrollmentRouteID int64
	Number            string
	OldNumber       *string
	IMEI            *string
	Phone           *string
	LastUpdate      int64
	Custom1         *string
	Custom2         *string
	Custom3         *string
	Info            *string
}
