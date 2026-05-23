package crypto

import "strings"

// NormalizeLoginPassword prepares the password field for PasswordMatch.
//
// Legacy clients (React/Angular/Java API) send MD5(rawPassword) as 32-char uppercase hex.
// Tools such as Swagger often send the raw password; we hash it the same way the UI does
// before compare, matching frontend loginPasswordEncode when publicKey is absent.
func NormalizeLoginPassword(password string) string {
	password = strings.TrimSpace(password)
	if IsMD5Hex(password) {
		return strings.ToUpper(password)
	}
	return MD5UpperHex(password)
}

// IsMD5Hex reports whether s looks like an MD5 digest in hex (32 characters).
func IsMD5Hex(s string) bool {
	if len(s) != 32 {
		return false
	}
	for _, c := range s {
		switch {
		case c >= '0' && c <= '9', c >= 'a' && c <= 'f', c >= 'A' && c <= 'F':
		default:
			return false
		}
	}
	return true
}
