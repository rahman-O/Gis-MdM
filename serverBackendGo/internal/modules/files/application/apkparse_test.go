package application

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestParseXAPKZip_fromManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bundle.zip")
	manifest := `{"package_name":"com.example.app","version_name":"1.2.3","version_code":42,"name":"Example App"}`
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	w, err := zw.Create("manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte(manifest)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	d := parseXAPKZip(path)
	if d == nil {
		t.Fatal("expected details")
	}
	if d.Pkg != "com.example.app" || d.Version != "1.2.3" || d.VersionCode != 42 || d.Name != "Example App" {
		t.Fatalf("unexpected details: %+v", d)
	}
}

func TestParseAPK_xapkRenamedToApk(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "myapp.apk")
	manifest := `{"package_name":"com.demo.pkg","version_name":"9.0","version_code":900,"name":"Demo"}`
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	w, err := zw.Create("manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte(manifest)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	d := ParseAPK(path)
	if d == nil || d.Pkg != "com.demo.pkg" || d.Version != "9.0" || d.VersionCode != 900 {
		t.Fatalf("expected XAPK manifest parse, got %+v", d)
	}
}
