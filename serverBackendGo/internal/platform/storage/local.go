package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const tempDelimiter = "1111111"

// LocalStore manages tenant-scoped files under a base directory.
type LocalStore struct {
	BaseDir string
}

func NewLocalStore(baseDir string) *LocalStore {
	return &LocalStore{BaseDir: filepath.Clean(baseDir)}
}

// IsSafePath rejects path traversal (Java FileUtil.isSafePath).
func IsSafePath(path string) bool {
	return path == "" || !strings.Contains(path, "..")
}

// AdjustFileName mirrors Java FileUtil.adjustFileName.
func AdjustFileName(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "+", "_")
	name = strings.ReplaceAll(name, "%", "_")
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")
	return name
}

// CreateTemp writes content to a temp file; returns absolute path.
func (s *LocalStore) CreateTemp(originalName string, r io.Reader) (string, error) {
	safe := AdjustFileName(originalName)
	tmp, err := os.CreateTemp("", safe+tempDelimiter+"*.temp")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	if _, err := io.Copy(tmp, r); err != nil {
		tmp.Close()
		os.Remove(path)
		return "", err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(path)
		return "", err
	}
	return path, nil
}

// NameFromTmpPath extracts original filename from Java-style temp name.
func NameFromTmpPath(tmpPath string) (string, error) {
	base := filepath.Base(tmpPath)
	if !strings.Contains(base, tempDelimiter) {
		return "", fmt.Errorf("invalid temp file name")
	}
	return strings.SplitN(base, tempDelimiter, 2)[0], nil
}

// CustomerRoot returns {BaseDir}/{filesDir}.
func (s *LocalStore) CustomerRoot(filesDir string) string {
	if filesDir == "" {
		return s.BaseDir
	}
	return filepath.Join(s.BaseDir, filesDir)
}

// MoveToCustomer moves temp file into tenant tree. relativePath is final filepath (may include subdirs).
func (s *LocalStore) MoveToCustomer(filesDir, subdir, tmpPath, fileName string) (string, error) {
	if !IsSafePath(subdir) || !IsSafePath(fileName) {
		return "", fmt.Errorf("unsafe path")
	}
	name := fileName
	if name == "" {
		var err error
		name, err = NameFromTmpPath(tmpPath)
		if err != nil {
			return "", err
		}
	}
	for strings.HasPrefix(name, "/") {
		name = name[1:]
	}
	rel := name
	if subdir != "" {
		rel = filepath.ToSlash(filepath.Join(subdir, name))
	}
	dest := filepath.Join(s.CustomerRoot(filesDir), filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", err
	}
	if _, err := os.Stat(dest); err == nil {
		return "", ErrExists
	} else if !os.IsNotExist(err) {
		return "", err
	}
	if err := os.Rename(tmpPath, dest); err != nil {
		if err := copyFile(tmpPath, dest); err != nil {
			return "", err
		}
		_ = os.Remove(tmpPath)
	}
	return filepath.ToSlash(rel), nil
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// DeleteRelative removes a file under customer root.
func (s *LocalStore) DeleteRelative(filesDir, relPath string) error {
	if !IsSafePath(relPath) {
		return fmt.Errorf("unsafe path")
	}
	p := filepath.Join(s.CustomerRoot(filesDir), filepath.FromSlash(relPath))
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(p)
}

// DirSizeBytes returns total size of customer directory.
func (s *LocalStore) DirSizeBytes(filesDir string) (int64, error) {
	root := s.CustomerRoot(filesDir)
	var total int64
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	return total, err
}

// BuildPublicURL builds {baseURL}/files/{filesDir}/{relPath}.
func BuildPublicURL(baseURL, filesDir, relPath string) string {
	base := strings.TrimRight(baseURL, "/")
	rel := strings.TrimLeft(strings.ReplaceAll(relPath, "\\", "/"), "/")
	if filesDir != "" {
		return fmt.Sprintf("%s/files/%s/%s", base, filesDir, rel)
	}
	return fmt.Sprintf("%s/files/%s", base, rel)
}

// ErrExists is returned when target file already exists.
var ErrExists = fmt.Errorf("file exists")
