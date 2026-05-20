package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsSafePath(t *testing.T) {
	if !IsSafePath("foo/bar.apk") {
		t.Fatal("expected safe")
	}
	if IsSafePath("../etc/passwd") {
		t.Fatal("expected unsafe")
	}
}

func TestCreateTempAndMove(t *testing.T) {
	dir := t.TempDir()
	store := NewLocalStore(dir)
	path, err := store.CreateTemp("my app.apk", strings.NewReader("apk-bytes"))
	if err != nil {
		t.Fatal(err)
	}
	rel, err := store.MoveToCustomer("customer-1", "", path, "my_app.apk")
	if err != nil {
		t.Fatal(err)
	}
	if rel != "my_app.apk" {
		t.Fatalf("rel=%q", rel)
	}
	full := filepath.Join(dir, "customer-1", "my_app.apk")
	if _, err := os.Stat(full); err != nil {
		t.Fatal(err)
	}
}

func TestMoveExists(t *testing.T) {
	dir := t.TempDir()
	store := NewLocalStore(dir)
	_ = os.MkdirAll(filepath.Join(dir, "c1"), 0o755)
	_ = os.WriteFile(filepath.Join(dir, "c1", "dup.txt"), []byte("x"), 0o644)
	tmp, _ := store.CreateTemp("dup.txt", strings.NewReader("y"))
	_, err := store.MoveToCustomer("c1", "", tmp, "dup.txt")
	if err != ErrExists {
		t.Fatalf("want ErrExists got %v", err)
	}
}
