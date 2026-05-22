package application

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/applications/domain"
	"github.com/gis-mdm/server-backend-go/internal/platform/storage"
)

func TestCommitApplicationAPK_setsURLFromTempUpload(t *testing.T) {
	base := t.TempDir()
	filesDir := "customer-1"
	customerRoot := filepath.Join(base, filesDir)
	if err := os.MkdirAll(customerRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	store := storage.NewLocalStore(base)
	tmpPath, err := store.CreateTemp("launcher.apk", &stringsReader{data: []byte("PK\x03\x04fake")})
	if err != nil {
		t.Fatal(err)
	}
	svc := &Service{
		repo:    appRepoStub{},
		store:   store,
		baseURL: "http://192.168.0.10:8080",
	}
	app := domain.Application{FilePath: strPtr(tmpPath)}
	if err := svc.commitApplicationAPK(context.Background(), 1, &app); err != nil {
		t.Fatal(err)
	}
	if app.URL == nil || *app.URL == "" {
		t.Fatal("expected url after commit")
	}
	if app.FilePath == nil || storage.IsTempUploadPath(*app.FilePath) {
		t.Fatalf("expected committed relative path, got %v", app.FilePath)
	}
}

type stringsReader struct {
	data []byte
	done bool
}

func (s *stringsReader) Read(p []byte) (int, error) {
	if s.done {
		return 0, io.EOF
	}
	n := copy(p, s.data)
	s.done = true
	return n, nil
}

func strPtr(s string) *string { return &s }
