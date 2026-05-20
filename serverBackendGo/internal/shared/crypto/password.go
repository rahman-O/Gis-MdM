package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
)

const passSalt = "5YdSYHyg2U"

// PASS_CHARS matches PasswordUtil.PASS_CHARS (subset used for auth tokens).
const passChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-.,!#$%()=+;*/"

const alphaCapsEnd = 61

func init() {
	rand.Seed(time.Now().UnixNano())
}

// HashFromMd5 computes SHA1(md5Hex + salt) as uppercase hex (legacy Headwind).
func HashFromMd5(md5Hex string) string {
	sum := sha1.Sum([]byte(md5Hex + passSalt))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

// PasswordMatch compares client MD5 (uppercase hex) with DB password hash.
func PasswordMatch(enteredMD5, dbPassword string) bool {
	return strings.EqualFold(HashFromMd5(enteredMD5), dbPassword)
}

// MD5UpperHex returns MD5 of raw password as uppercase hex (for tests).
func MD5UpperHex(raw string) string {
	sum := md5.Sum([]byte(raw))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

// GenerateAuthToken creates a 20-character token like PasswordUtil.generateToken().
func GenerateAuthToken() string {
	b := make([]byte, 20)
	for i := range b {
		b[i] = passChars[rand.Intn(alphaCapsEnd+1)]
	}
	return string(b)
}
