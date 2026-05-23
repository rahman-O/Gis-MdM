package crypto

import "testing"

func TestPasswordMatch_knownVector(t *testing.T) {
	md5 := MD5UpperHex("admin")
	hash := HashFromMd5(md5)
	if !PasswordMatch(md5, hash) {
		t.Fatalf("expected match for derived hash")
	}
	if PasswordMatch(md5, "wrong") {
		t.Fatal("expected no match")
	}
}

func TestGenerateAuthToken_length(t *testing.T) {
	tok := GenerateAuthToken()
	if len(tok) != 20 {
		t.Fatalf("expected len 20, got %d", len(tok))
	}
}
