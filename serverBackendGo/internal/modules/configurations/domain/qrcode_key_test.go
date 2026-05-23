package domain

import (
	"strings"
	"testing"
)

func TestNewQRCodeKey_format(t *testing.T) {
	k := NewQRCodeKey()
	if len(k) != 32 {
		t.Fatalf("len=%d want 32", len(k))
	}
	for _, c := range k {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			t.Fatalf("non-hex char %q in %q", c, k)
		}
	}
}

func TestEnsureQRCodeKey_preservesPayload(t *testing.T) {
	existing := "abc"
	cfg := Configuration{QRCodeKey: &existing}
	EnsureQRCodeKey(&cfg, nil)
	if cfg.QRCodeKey == nil || *cfg.QRCodeKey != existing {
		t.Fatalf("got %v", cfg.QRCodeKey)
	}
}

func TestEnsureQRCodeKey_generatesWhenEmpty(t *testing.T) {
	cfg := Configuration{}
	EnsureQRCodeKey(&cfg, nil)
	if cfg.QRCodeKey == nil || strings.TrimSpace(*cfg.QRCodeKey) == "" {
		t.Fatal("expected generated key")
	}
}

func TestEnsureQRCodeKey_preservesDatabase(t *testing.T) {
	dbKey := "db-key"
	db := Configuration{QRCodeKey: &dbKey}
	cfg := Configuration{}
	EnsureQRCodeKey(&cfg, &db)
	if cfg.QRCodeKey == nil || *cfg.QRCodeKey != dbKey {
		t.Fatalf("got %v", cfg.QRCodeKey)
	}
}
