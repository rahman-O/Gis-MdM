package domain

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"strings"
)

// NewQRCodeKey returns a unique enrollment key (32-char hex, same style as Java MD5(RANDOM())).
func NewQRCodeKey() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	sum := md5.Sum(b[:])
	return hex.EncodeToString(sum[:])
}

// EnsureQRCodeKey assigns a new key when missing; preserves existingDB when the payload omits a key.
func EnsureQRCodeKey(cfg *Configuration, existingDB *Configuration) {
	if cfg == nil {
		return
	}
	if cfg.QRCodeKey != nil && strings.TrimSpace(*cfg.QRCodeKey) != "" {
		return
	}
	if existingDB != nil && existingDB.QRCodeKey != nil && strings.TrimSpace(*existingDB.QRCodeKey) != "" {
		cfg.QRCodeKey = existingDB.QRCodeKey
		return
	}
	k := NewQRCodeKey()
	cfg.QRCodeKey = &k
}
