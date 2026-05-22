package storage

import "strings"

// IsTempUploadPath reports Java/Go temp APK paths (CreateTemp uses tempDelimiter in the basename).
func IsTempUploadPath(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}
	return strings.Contains(path, tempDelimiter) || strings.HasSuffix(strings.ToLower(path), ".temp")
}
