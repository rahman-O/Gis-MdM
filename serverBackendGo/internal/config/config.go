package config

import (
	"os"
	"strconv"
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

	FilesDirectory string
	HashSecret     string

	RebrandingName      string
	RebrandingLogo      string
	RebrandingVendor    string
	RebrandingVendorURL string
	RebrandingSignupURL string
	RebrandingTermsURL  string
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

		FilesDirectory: getenv("FILES_DIRECTORY", "/var/lib/hmdm/files"),
		HashSecret:     getenv("HASH_SECRET", "changeme-C3z9vi54"),

		RebrandingName:      getenv("REBRANDING_NAME", "Headwind MDM"),
		RebrandingLogo:      getenv("REBRANDING_LOGO", ""),
		RebrandingVendor:    getenv("REBRANDING_VENDOR_NAME", ""),
		RebrandingVendorURL: getenv("REBRANDING_VENDOR_LINK", ""),
		RebrandingSignupURL: getenv("REBRANDING_SIGNUP_LINK", ""),
		RebrandingTermsURL:  getenv("REBRANDING_TERMS_LINK", ""),
	}
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
