package storage

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestServeFileRel_servesFile(t *testing.T) {
	base := t.TempDir()
	name := "hmdm-test.apk"
	if err := os.WriteFile(filepath.Join(base, name), []byte("PK\x03\x04"), 0o644); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/files/"+name, nil)
	rec := httptest.NewRecorder()
	ServeFileRel(rec, req, base, name)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}
