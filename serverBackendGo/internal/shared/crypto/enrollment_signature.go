package crypto

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"strings"
)

var wsStrip = regexp.MustCompile(`\s+`)

// SHA1UpperHex returns SHA1 digest as uppercase hex (legacy CryptoUtil.getSHA1String).
func SHA1UpperHex(raw string) string {
	sum := sha1.Sum([]byte(raw))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

// CheckRequestSignature validates X-Request-Signature against SHA1(value).
func CheckRequestSignature(signature, value string) bool {
	if strings.TrimSpace(signature) == "" {
		return false
	}
	want := SHA1UpperHex(value)
	return strings.EqualFold(strings.TrimSpace(signature), want)
}

// SignSyncResponse returns SHA1(hashSecret + compactJSON) for X-Response-Signature.
func SignSyncResponse(hashSecret string, payload any) string {
	b, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	compact := wsStrip.ReplaceAllString(string(b), "")
	return SHA1UpperHex(hashSecret + compact)
}
