package application

import (
	"archive/zip"
	"path/filepath"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/files/domain"
)

// ParseAPK extracts best-effort metadata from an APK on disk.
func ParseAPK(path string) *domain.APKFileDetails {
	if !strings.HasSuffix(strings.ToLower(path), ".apk") {
		return nil
	}
	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil
	}
	defer zr.Close()
	details := &domain.APKFileDetails{}
	base := filepath.Base(path)
	if strings.Contains(base, "arm64") || strings.Contains(base, "arm64-v8a") {
		details.Arch = "arm64"
	} else if strings.Contains(base, "armeabi") {
		details.Arch = "armeabi"
	}
	for _, f := range zr.File {
		name := strings.ToLower(f.Name)
		if strings.HasSuffix(name, ".apk") && name != "androidmanifest.xml" {
			// split APK bundle — keep arch hint from nested name
			if strings.Contains(name, "arm64") {
				details.Arch = "arm64"
			}
		}
	}
	// Without a full binary XML parser, derive package hint from filename when possible.
	stem := strings.TrimSuffix(base, ".apk")
	if details.Pkg == "" && strings.Contains(stem, ".") {
		details.Pkg = stem
	}
	return details
}
