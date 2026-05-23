package application

import (
	"archive/zip"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/files/domain"
	"github.com/gis-mdm/server-backend-go/internal/platform/storage"
	"github.com/shogo82148/androidbinary/apk"
)

// ParseAPK extracts metadata from an APK or XAPK on disk (parity with Java APKFileAnalyzer).
func ParseAPK(path string) *domain.APKFileDetails {
	if path == "" {
		return nil
	}
	if d := parseXAPKZip(path); d != nil && d.Pkg != "" {
		return d
	}
	return parseBinaryAPK(path)
}

func parseBinaryAPK(path string) *domain.APKFileDetails {
	lower := strings.ToLower(path)
	if !strings.HasSuffix(lower, ".apk") && !strings.HasSuffix(lower, ".temp") {
		return nil
	}
	f, err := apk.OpenFile(path)
	if err != nil {
		return fallbackDetailsFromFileName(path)
	}
	defer f.Close()

	manifest := f.Manifest()
	details := &domain.APKFileDetails{
		Pkg: manifest.Package.MustString(),
	}
	if v, err := manifest.VersionName.String(); err == nil {
		details.Version = v
	}
	details.VersionCode = int(manifest.VersionCode.MustInt32())
	if label, err := f.Label(nil); err == nil && strings.TrimSpace(label) != "" {
		details.Name = strings.TrimSpace(label)
	}
	details.Arch = archFromNativeLibs(path)
	if details.Pkg == "" {
		return fallbackDetailsFromFileName(path)
	}
	return details
}

type xapkManifest struct {
	PackageName string `json:"package_name"`
	VersionName string `json:"version_name"`
	VersionCode int    `json:"version_code"`
	Name        string `json:"name"`
}

func parseXAPKZip(path string) *domain.APKFileDetails {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil
	}
	defer zr.Close()
	var manifestRaw string
	for _, f := range zr.File {
		if strings.EqualFold(f.Name, "manifest.json") {
			rc, err := f.Open()
			if err != nil {
				return nil
			}
			b, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil
			}
			manifestRaw = string(b)
			break
		}
	}
	if manifestRaw == "" {
		return nil
	}
	var m xapkManifest
	if err := json.Unmarshal([]byte(manifestRaw), &m); err != nil {
		return nil
	}
	details := &domain.APKFileDetails{
		Pkg:         strings.TrimSpace(m.PackageName),
		Name:        strings.TrimSpace(m.Name),
		Version:     strings.TrimSpace(m.VersionName),
		VersionCode: m.VersionCode,
	}
	hasArm64 := strings.Contains(manifestRaw, "arm64")
	hasArmeabi := strings.Contains(manifestRaw, "armeabi")
	if hasArm64 && !hasArmeabi {
		details.Arch = "arm64"
	} else if hasArmeabi && !hasArm64 {
		details.Arch = "armeabi"
	}
	return details
}

// archFromNativeLibs mirrors Java getArchByApkLibs (armeabi-v7a / arm64-v8a).
func archFromNativeLibs(path string) string {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return archHintFromFileName(path)
	}
	defer zr.Close()
	abis := make(map[string]struct{})
	for _, f := range zr.File {
		if !strings.HasPrefix(f.Name, "lib/") {
			continue
		}
		parts := strings.Split(f.Name, "/")
		if len(parts) > 1 && parts[1] != "" {
			abis[parts[1]] = struct{}{}
		}
	}
	var result string
	for abi := range abis {
		switch abi {
		case "arm64-v8a":
			if result == "armeabi" {
				return ""
			}
			result = "arm64"
		case "armeabi-v7a":
			if result == "arm64" {
				return ""
			}
			result = "armeabi"
		}
	}
	if result != "" {
		return result
	}
	return archHintFromFileName(path)
}

func archHintFromFileName(path string) string {
	base := strings.ToLower(filepath.Base(path))
	if strings.Contains(base, "arm64") || strings.Contains(base, "arm64-v8a") {
		return "arm64"
	}
	if strings.Contains(base, "armeabi") {
		return "armeabi"
	}
	return ""
}

func fallbackDetailsFromFileName(path string) *domain.APKFileDetails {
	name, err := storage.NameFromTmpPath(path)
	if err != nil {
		name = filepath.Base(path)
	}
	stem := strings.TrimSuffix(name, ".apk")
	stem = strings.TrimSuffix(stem, ".xapk")
	details := &domain.APKFileDetails{Arch: archHintFromFileName(path)}
	if strings.Contains(stem, ".") {
		details.Pkg = stem
	}
	return details
}
