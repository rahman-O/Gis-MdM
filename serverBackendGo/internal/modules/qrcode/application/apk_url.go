package application

import (
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/port"
	"github.com/gis-mdm/server-backend-go/internal/platform/storage"
)

// ResolveVersionDownloadURL is the main-app version URL used for checksum (Java uses appVersion.getUrl()).
func ResolveVersionDownloadURL(cfg *port.QRConfig, baseURL string) string {
	return resolveDownloadURL(cfg, baseURL, false)
}

// ResolveMainAppDownloadURL picks launcher override, then version / app / filepath URLs.
func ResolveMainAppDownloadURL(cfg *port.QRConfig, baseURL string) string {
	return resolveDownloadURL(cfg, baseURL, true)
}

func resolveDownloadURL(cfg *port.QRConfig, baseURL string, allowLauncher bool) string {
	if cfg == nil {
		return ""
	}
	if allowLauncher {
		if u := strings.TrimSpace(cfg.LauncherURL); u != "" {
			return u
		}
	}
	if u := strings.TrimSpace(cfg.MainAppURL); u != "" {
		return u
	}
	if u := strings.TrimSpace(cfg.AppLevelURL); u != "" {
		return u
	}
	fp := strings.TrimSpace(cfg.MainAppFilePath)
	filesDir := strings.TrimSpace(cfg.FilesDir)
	if fp != "" && filesDir != "" && strings.TrimSpace(baseURL) != "" {
		return storage.BuildPublicURL(baseURL, filesDir, fp)
	}
	return ""
}
