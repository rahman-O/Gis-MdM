package domain

import "testing"

func TestNormalizePolicyLocks_filtersUnknown(t *testing.T) {
	got := NormalizePolicyLocks(map[string]bool{
		"mainAppId": true,
		"unknown":   true,
		"kioskMode": false,
	})
	if len(got) != 1 || !got["mainAppId"] {
		t.Fatalf("got %#v", got)
	}
}

func TestApplicationSettingLockKey(t *testing.T) {
	k := ApplicationSettingLockKey("com.hmdm.launcher", "debug")
	if k != "applicationSetting.com.hmdm.launcher.debug" {
		t.Fatal(k)
	}
}
