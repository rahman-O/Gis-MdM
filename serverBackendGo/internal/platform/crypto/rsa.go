package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const (
	privateKeyFile = "private.key"
	publicKeyFile  = "public.key"
)

// RSAKeys manages RSA key pair for transmit.password (legacy Headwind).
type RSAKeys struct {
	dir string
	mu  sync.Mutex
}

// NewRSAKeys stores keys under dir (e.g. FILES_DIRECTORY).
func NewRSAKeys(dir string) *RSAKeys {
	return &RSAKeys{dir: dir}
}

// EnsureKeys generates keys if missing.
func (k *RSAKeys) EnsureKeys() error {
	k.mu.Lock()
	defer k.mu.Unlock()
	if err := os.MkdirAll(k.dir, 0o700); err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(k.dir, publicKeyFile)); err == nil {
		return nil
	}
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}
	pubDER, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return err
	}
	privDER := x509.MarshalPKCS1PrivateKey(priv)
	if err := os.WriteFile(filepath.Join(k.dir, privateKeyFile), privDER, 0o600); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(k.dir, publicKeyFile), pubDER, 0o644)
}

// PublicKeyBase64 returns PKIX public key bytes as base64 (for auth options).
func (k *RSAKeys) PublicKeyBase64() (string, error) {
	if err := k.EnsureKeys(); err != nil {
		return "", err
	}
	b, err := os.ReadFile(filepath.Join(k.dir, publicKeyFile))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// DecryptBase64 decrypts RSA/PKCS1 ciphertext from client (base64).
func (k *RSAKeys) DecryptBase64(cipherB64 string) (string, error) {
	if err := k.EnsureKeys(); err != nil {
		return "", err
	}
	raw, err := base64.StdEncoding.DecodeString(cipherB64)
	if err != nil {
		return "", err
	}
	privDER, err := os.ReadFile(filepath.Join(k.dir, privateKeyFile))
	if err != nil {
		return "", err
	}
	priv, err := x509.ParsePKCS1PrivateKey(privDER)
	if err != nil {
		return "", err
	}
	plain, err := rsa.DecryptPKCS1v15(rand.Reader, priv, raw)
	if err != nil {
		return "", fmt.Errorf("rsa decrypt: %w", err)
	}
	return string(plain), nil
}
