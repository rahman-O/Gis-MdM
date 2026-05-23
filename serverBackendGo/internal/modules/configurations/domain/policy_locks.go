package domain

import "strings"

const policyLocksKey = "policyLocks"

// PolicyLocksKey is the settingsjson property for field lock map.
func PolicyLocksKey() string { return policyLocksKey }

// AllowedPolicyLockFields are MDM editor fields that may be locked at configuration level.
var AllowedPolicyLockFields = map[string]struct{}{
	"mainAppId": {}, "contentAppId": {}, "kioskMode": {}, "restrictions": {},
	"gps": {}, "bluetooth": {}, "wifi": {}, "mobileData": {}, "usbStorage": {},
	"lockSafeSettings": {}, "lockVolume": {}, "passwordMode": {},
	"eventReceivingComponent": {}, "launcherUrl": {},
}

// NormalizePolicyLocks keeps only allowed keys with true values.
func NormalizePolicyLocks(in map[string]bool) map[string]bool {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]bool)
	for k, v := range in {
		k = strings.TrimSpace(k)
		if k == "" || !v {
			continue
		}
		if _, ok := AllowedPolicyLockFields[k]; ok {
			out[k] = true
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// ApplicationSettingLockKey builds a policyLocks entry for per-app settings.
func ApplicationSettingLockKey(pkg, name string) string {
	return "applicationSetting." + strings.TrimSpace(pkg) + "." + strings.TrimSpace(name)
}

// IsApplicationSettingLockKey reports keys created by ApplicationSettingLockKey.
func IsApplicationSettingLockKey(key string) bool {
	return strings.HasPrefix(key, "applicationSetting.")
}
