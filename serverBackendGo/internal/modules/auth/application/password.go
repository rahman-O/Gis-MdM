package application

import (
	sharedcrypto "github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

// ResolvePassword normalizes client password (MD5 hex or raw; optional RSA ciphertext).
func (s *Service) ResolvePassword(password string) (string, error) {
	if s.transmitPassword && s.rsa != nil && !sharedcrypto.IsMD5Hex(password) {
		raw, err := s.rsa.DecryptBase64(password)
		if err != nil {
			return "", err
		}
		return sharedcrypto.NormalizeLoginPassword(raw), nil
	}
	return sharedcrypto.NormalizeLoginPassword(password), nil
}
