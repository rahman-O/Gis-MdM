package storage

import "testing"

func TestSafeFilePath(t *testing.T) {
	base := t.TempDir()
	okPath, ok := SafeFilePath(base, "customer-1/apk/app.apk")
	if !ok || okPath == "" {
		t.Fatal("expected valid path")
	}
	if _, ok := SafeFilePath(base, "../etc/passwd"); ok {
		t.Fatal("expected traversal rejected")
	}
	if _, ok := SafeFilePath(base, ""); ok {
		t.Fatal("expected empty rejected")
	}
}
