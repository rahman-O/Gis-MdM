package storage

import (
	"net/http"
	"path/filepath"
	"strings"
)

// SafeFilePath resolves rel under baseDir and rejects path traversal.
func SafeFilePath(baseDir, rel string) (string, bool) {
	raw := strings.TrimSpace(rel)
	if raw == "" || strings.Contains(raw, "..") {
		return "", false
	}
	rel = strings.TrimPrefix(filepath.Clean("/"+raw), "/")
	if rel == "" || rel == "." {
		return "", false
	}
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", false
	}
	full := filepath.Join(absBase, filepath.FromSlash(rel))
	absFull, err := filepath.Abs(full)
	if err != nil {
		return "", false
	}
	if absFull != absBase && !strings.HasPrefix(absFull, absBase+string(filepath.Separator)) {
		return "", false
	}
	relToBase, err := filepath.Rel(absBase, absFull)
	if err != nil || strings.HasPrefix(relToBase, "..") {
		return "", false
	}
	return absFull, true
}

// ServeFileRel serves one file under baseDir (relative path within FILES_DIRECTORY).
func ServeFileRel(w http.ResponseWriter, r *http.Request, baseDir, rel string) {
	baseDir = strings.TrimSpace(baseDir)
	if baseDir == "" {
		http.NotFound(w, r)
		return
	}
	rel = strings.TrimPrefix(strings.TrimSpace(rel), "/")
	full, ok := SafeFilePath(baseDir, rel)
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, full)
}

// ServeFilesHandler serves GET /files/{customerDir}/... from baseDir (FILES_DIRECTORY).
func ServeFilesHandler(baseDir string) http.HandlerFunc {
	baseDir = strings.TrimSpace(baseDir)
	return func(w http.ResponseWriter, r *http.Request) {
		rel := strings.TrimPrefix(r.URL.Path, "/files/")
		ServeFileRel(w, r, baseDir, rel)
	}
}
