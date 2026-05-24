package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds environment-driven settings.
type Config struct {
	ServerPort string
	GinMode    string

	DatabaseURL string

	JWTSecret               string
	JWTExpiryHours          int
	JWTValiditySeconds      int64
	JWTValidityRememberSecs int64

	SessionSecret string

	CustomerSignup   bool
	EmailConfigured  bool
	TransmitPassword bool
	BaseURL          string
	SwaggerEnabled   bool

	PublicIPAllowlist  string
	PrivateIPAllowlist string

	ModuleAuthEnabled          bool
	ModuleSignupEnabled        bool
	ModulePasswordResetEnabled bool
	ModuleFilesEnabled         bool
	ModuleIconsEnabled         bool
	ModulePublicAPIEnabled     bool
	ModuleSyncEnabled          bool
	ModulePushEnabled          bool
	ModuleNotificationsEnabled bool
	ModuleUpdatesEnabled       bool
	ModuleQRCodeEnabled        bool
	ModuleDeviceTreeEnabled    bool
	ModuleProfilesEnabled      bool
	ModuleProfileRolloutEnabled bool
	ModuleEnrollmentRoutesEnabled bool
	ProfileStalePublishDays       int

	FilesDirectory string
	HashSecret     string

	SecureEnrollment           bool
	PreventDuplicateEnrollment bool
	PollingTimeoutMs           int

	RebrandingName       string
	RebrandingLogo       string
	RebrandingVendor     string
	RebrandingVendorURL  string
	RebrandingSignupURL  string
	RebrandingTermsURL   string
	RebrandingMobileName string

	UpdateManifestURL string

	EnabledPlugins              []string
	ModulePluginsEnabled        bool
	ModulePluginsPlatformEnabled bool
	ModulePluginsAuditEnabled   bool
	ModulePluginsMessagingEnabled bool
	ModulePluginsDeviceinfoEnabled bool
	ModulePluginsDevicelogEnabled bool

	PushNotifierEnabled      bool
	PushScheduleIntervalSec  int
	ModuleStatsEnabled       bool
	ModuleVideosEnabled      bool
	VideoDirectory           string
}

// Load reads configuration from environment variables.
func Load() Config {
	return Config{
		ServerPort: getenv("SERVER_PORT", "8080"),
		GinMode:    getenv("GIN_MODE", "debug"),

		DatabaseURL: os.Getenv("DATABASE_URL"),

		JWTSecret:               getenv("JWT_SECRET", "change-me"),
		JWTExpiryHours:          getenvInt("JWT_EXPIRY_HOURS", 24),
		JWTValiditySeconds:      int64(getenvInt("JWT_VALIDITY_SECONDS", 86400)),
		JWTValidityRememberSecs: int64(getenvInt("JWT_VALIDITY_REMEMBER_SECONDS", 2592000)),

		SessionSecret:   getenv("SESSION_SECRET", "change-me-session"),
		CustomerSignup:   getenvBool("CUSTOMER_SIGNUP", false),
		EmailConfigured:  getenvBool("EMAIL_CONFIGURED", false),
		TransmitPassword: getenvBool("TRANSMIT_PASSWORD", false),
		BaseURL:          getenv("BASE_URL", "http://localhost:8080"),
		SwaggerEnabled:   getenvBool("SWAGGER_ENABLED", true),

		PublicIPAllowlist:  os.Getenv("PUBLIC_IP_ALLOWLIST"),
		PrivateIPAllowlist: os.Getenv("PRIVATE_IP_ALLOWLIST"),

		ModuleAuthEnabled:          getenvBool("MODULE_AUTH_ENABLED", true),
		ModuleSignupEnabled:        getenvBool("MODULE_SIGNUP_ENABLED", false),
		ModulePasswordResetEnabled: getenvBool("MODULE_PASSWORDRESET_ENABLED", false),
		ModuleFilesEnabled:         getenvBool("MODULE_FILES_ENABLED", true),
		ModuleIconsEnabled:         getenvBool("MODULE_ICONS_ENABLED", true),
		ModulePublicAPIEnabled:     getenvBool("MODULE_PUBLICAPI_ENABLED", true),
		ModuleSyncEnabled:          getenvBool("MODULE_SYNC_ENABLED", true),
		ModulePushEnabled:          getenvBool("MODULE_PUSH_ENABLED", true),
		ModuleNotificationsEnabled: getenvBool("MODULE_NOTIFICATIONS_ENABLED", true),
		ModuleUpdatesEnabled:       getenvBool("MODULE_UPDATES_ENABLED", true),
		ModuleQRCodeEnabled:        getenvBool("MODULE_QRCODE_ENABLED", true),
		ModuleDeviceTreeEnabled:    getenvBool("MODULE_DEVICE_TREE_ENABLED", true),
		ModuleProfilesEnabled:      getenvBool("MODULE_PROFILES_ENABLED", true),
		ModuleProfileRolloutEnabled: getenvBool("MODULE_PROFILE_ROLLOUT_ENABLED", true),
		ModuleEnrollmentRoutesEnabled: getenvBool("MODULE_ENROLLMENT_ROUTES_ENABLED", true),
		ProfileStalePublishDays:       getenvInt("PROFILE_STALE_PUBLISH_DAYS", 30),

		FilesDirectory: getenv("FILES_DIRECTORY", "/var/lib/hmdm/files"),
		HashSecret:     getenv("HASH_SECRET", "changeme-C3z9vi54"),

		SecureEnrollment:           getenvBool("SECURE_ENROLLMENT", false),
		PreventDuplicateEnrollment: getenvBool("PREVENT_DUPLICATE_ENROLLMENT", false),
		PollingTimeoutMs:           getenvInt("POLLING_TIMEOUT_MS", 60000),

		RebrandingName:       getenv("REBRANDING_NAME", "Headwind MDM"),
		RebrandingLogo:       getenv("REBRANDING_LOGO", ""),
		RebrandingVendor:     getenv("REBRANDING_VENDOR_NAME", ""),
		RebrandingVendorURL:  getenv("REBRANDING_VENDOR_LINK", ""),
		RebrandingSignupURL:  getenv("REBRANDING_SIGNUP_LINK", ""),
		RebrandingTermsURL:   getenv("REBRANDING_TERMS_LINK", ""),
		RebrandingMobileName: getenv("REBRANDING_MOBILE_NAME", ""),

		UpdateManifestURL: getenv("UPDATE_MANIFEST_URL", "https://h-mdm.com/files/hmdm_update_manifest.txt"),

		EnabledPlugins:               parseEnabledPlugins(getenv("ENABLED_PLUGINS", "audit,push,messaging,deviceinfo,devicelog")),
		ModulePluginsEnabled:         getenvBool("MODULE_PLUGINS_ENABLED", true),
		ModulePluginsPlatformEnabled: getenvBool("MODULE_PLUGINS_PLATFORM_ENABLED", true),
		ModulePluginsAuditEnabled:    getenvBool("MODULE_PLUGINS_AUDIT_ENABLED", true),
		ModulePluginsMessagingEnabled: getenvBool("MODULE_PLUGINS_MESSAGING_ENABLED", true),
		ModulePluginsDeviceinfoEnabled: getenvBool("MODULE_PLUGINS_DEVICEINFO_ENABLED", true),
		ModulePluginsDevicelogEnabled:  getenvBool("MODULE_PLUGINS_DEVICELOG_ENABLED", true),

		PushNotifierEnabled:     getenvBool("MODULE_PUSH_NOTIFIER_ENABLED", true),
		PushScheduleIntervalSec: getenvInt("PUSH_SCHEDULE_INTERVAL_SEC", 60),
		ModuleStatsEnabled:      getenvBool("MODULE_STATS_ENABLED", false),
		ModuleVideosEnabled:     getenvBool("MODULE_VIDEOS_ENABLED", false),
		VideoDirectory:          getenv("VIDEO_DIRECTORY", "./data/videos"),
	}
}

func parseEnabledPlugins(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(strings.ToLower(p))
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// IsPluginEnabled returns true if identifier is in ENABLED_PLUGINS list.
func (c Config) IsPluginEnabled(identifier string) bool {
	id := strings.ToLower(strings.TrimSpace(identifier))
	for _, e := range c.EnabledPlugins {
		if e == id {
			return true
		}
	}
	return false
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getenvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	switch v {
	case "1", "true", "TRUE", "yes", "YES":
		return true
	case "0", "false", "FALSE", "no", "NO":
		return false
	default:
		return fallback
	}
}
