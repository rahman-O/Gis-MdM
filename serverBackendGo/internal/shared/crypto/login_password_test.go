package crypto

import (
	"strings"
	"testing"
)

func TestNormalizeLoginPassword_raw(t *testing.T) {
	got := NormalizeLoginPassword("admin")
	want := MD5UpperHex("admin")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestNormalizeLoginPassword_alreadyMD5(t *testing.T) {
	md5 := MD5UpperHex("admin")
	if NormalizeLoginPassword(md5) != md5 {
		t.Fatal("expected unchanged MD5")
	}
	if NormalizeLoginPassword(strings.ToLower(md5)) != md5 {
		t.Fatal("expected uppercase MD5")
	}
}
