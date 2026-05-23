package crypto

import "math/rand"

// GeneratePassword creates a random password (legacy PasswordUtil.generatePassword, strength 0).
func GeneratePassword(length int) string {
	if length < 8 {
		length = 8
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = passChars[rand.Intn(alphaCapsEnd+1)]
	}
	return string(b)
}
