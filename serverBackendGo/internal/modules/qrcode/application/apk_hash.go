package application

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ApkChecksum returns SHA-256 digest as base64 (Java CryptoUtil.getBase64String).
func ApkChecksum(apkURL, baseURL, filesDirectory string, storedHash string) (string, error) {
	if strings.TrimSpace(storedHash) != "" {
		return strings.TrimSpace(storedHash), nil
	}
	apkURL = strings.TrimSpace(apkURL)
	if apkURL == "" {
		return "", fmt.Errorf("empty apk url")
	}
	if localPath := localFilesPath(apkURL, baseURL, filesDirectory); localPath != "" {
		return hashFile(localPath)
	}
	// Absolute temp upload path still on disk (legacy rows before APK commit fix).
	if strings.HasPrefix(apkURL, "/") {
		if st, err := os.Stat(apkURL); err == nil && !st.IsDir() {
			return hashFile(apkURL)
		}
	}
	return hashRemote(apkURL)
}

func localFilesPath(apkURL, baseURL, filesDirectory string) string {
	u, err := url.Parse(apkURL)
	if err != nil {
		return ""
	}
	path := u.Path
	idx := strings.Index(path, "/files/")
	if idx < 0 {
		return ""
	}
	rel := strings.TrimPrefix(path[idx+len("/files/"):], "/")
	if rel == "" {
		return ""
	}
	full, ok := safeJoin(filesDirectory, rel)
	if !ok {
		return ""
	}
	if _, err := os.Stat(full); err != nil {
		return ""
	}
	_ = baseURL
	return full
}

func safeJoin(baseDir, rel string) (string, bool) {
	rel = strings.TrimPrefix(filepath.Clean("/"+rel), "/")
	if rel == "" || strings.Contains(rel, "..") {
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
	return absFull, true
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return hashReader(f)
}

func hashRemote(apkURL string) (string, error) {
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Get(apkURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("apk download status %d", resp.StatusCode)
	}
	return hashReader(resp.Body)
}

func hashReader(r io.Reader) (string, error) {
	h := sha256.New()
	buf := make([]byte, 8192)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			h.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
