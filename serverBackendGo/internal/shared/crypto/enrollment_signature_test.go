package crypto

import "testing"

func TestCheckRequestSignature(t *testing.T) {
	secret := "changeme-C3z9vi54"
	deviceID := "hmdm-001"
	want := SHA1UpperHex(secret + deviceID)
	if !CheckRequestSignature(want, secret+deviceID) {
		t.Fatalf("expected valid signature %s", want)
	}
	if CheckRequestSignature("BAD", secret+deviceID) {
		t.Fatal("expected invalid signature")
	}
}

func TestSignSyncResponseStable(t *testing.T) {
	payload := map[string]any{"deviceId": "hmdm-001", "status": "ok"}
	s1 := SignSyncResponse("secret", payload)
	s2 := SignSyncResponse("secret", payload)
	if s1 == "" || s1 != s2 {
		t.Fatalf("signature mismatch: %q %q", s1, s2)
	}
}
