package crypto

// DeviceUploadHash returns MD5(deviceID + secret) uppercase hex (PublicResource / AppList).
func DeviceUploadHash(deviceID, secret string) string {
	return MD5UpperHex(deviceID + secret)
}
