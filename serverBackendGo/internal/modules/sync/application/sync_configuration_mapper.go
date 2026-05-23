package application

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/sync/domain"
)

// ApplyConfigurationPolicy maps settingsjson and column extras onto SyncResponse (Java SyncResource parity).
func ApplyConfigurationPolicy(resp *domain.SyncResponse, settingsJSON []byte, backgroundImageURL *string) {
	if resp == nil {
		return
	}
	if backgroundImageURL != nil && *backgroundImageURL != "" {
		resp.BackgroundImageURL = backgroundImageURL
	}
	if len(settingsJSON) == 0 {
		return
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(settingsJSON, &m); err != nil {
		return
	}
	resp.GPS = boolPtr(m, "gps")
	resp.Bluetooth = boolPtr(m, "bluetooth")
	resp.Wifi = boolPtr(m, "wifi")
	resp.MobileData = boolPtr(m, "mobileData")
	resp.UsbStorage = boolPtr(m, "usbStorage")
	resp.LockStatusBar = boolPtr(m, "blockStatusBar")
	if resp.LockStatusBar == nil {
		resp.LockStatusBar = boolPtr(m, "lockStatusBar")
	}
	resp.DisableLocation = boolPtr(m, "disableLocation")
	resp.KioskMode = boolVal(m, "kioskMode")
	resp.Restrictions = stringPtr(m, "restrictions")
	resp.KioskHome = boolPtrTrueOnly(m, "kioskHome")
	resp.KioskRecents = boolPtrTrueOnly(m, "kioskRecents")
	resp.KioskNotifications = boolPtrTrueOnly(m, "kioskNotifications")
	resp.KioskSystemInfo = boolPtrTrueOnly(m, "kioskSystemInfo")
	resp.KioskKeyguard = boolPtrTrueOnly(m, "kioskKeyguard")
	resp.KioskLockButtons = boolPtrTrueOnly(m, "kioskLockButtons")
	resp.KioskScreenOn = boolPtrTrueOnly(m, "kioskScreenOn")
	resp.KioskExit = stringPtr(m, "kioskExit")
	resp.ShowWifi = boolPtr(m, "showWifi")
	resp.LockSafeSettings = boolPtr(m, "lockSafeSettings")
	resp.LockVolume = boolPtr(m, "lockVolume")
	resp.AutoBrightness = boolPtr(m, "autoBrightness")
	resp.ManageTimeout = boolPtrTrueOnly(m, "manageTimeout")
	resp.ManageVolume = boolPtrTrueOnly(m, "manageVolume")
	resp.RunDefaultLauncher = boolPtr(m, "runDefaultLauncher")
	resp.DisableScreenshots = boolPtr(m, "disableScreenshots")
	resp.AutostartForeground = boolPtr(m, "autostartForeground")
	resp.ScheduleAppUpdate = boolVal(m, "scheduleAppUpdate")
	resp.AppPermissions = stringPtr(m, "appPermissions")
	resp.PasswordMode = stringPtr(m, "passwordMode")
	resp.TimeZone = stringPtr(m, "timeZone")
	resp.AllowedClasses = stringPtr(m, "allowedClasses")
	resp.NewServerUrl = stringPtr(m, "newServerUrl")
	resp.MainApp = stringPtr(m, "mainApp")
	resp.SystemUpdateType = intPtr(m, "systemUpdateType")
	resp.Orientation = intPtr(m, "orientation")
	resp.Brightness = intPtr(m, "brightness")
	resp.Timeout = intPtr(m, "timeout")
	resp.Volume = intPtr(m, "volume")
	resp.KeepaliveTime = intPtr(m, "keepaliveTime")
	resp.IconSize = iconSizePtr(m, "iconSize")
	resp.SystemUpdateFrom = stringPtr(m, "systemUpdateFrom")
	resp.SystemUpdateTo = stringPtr(m, "systemUpdateTo")
	resp.AppUpdateFrom = stringPtr(m, "appUpdateFrom")
	resp.AppUpdateTo = stringPtr(m, "appUpdateTo")
	resp.DownloadUpdates = stringPtr(m, "downloadUpdates")
	if v := stringPtr(m, "pushOptions"); v != nil {
		resp.PushOptions = v
	}
	if v := stringPtr(m, "requestUpdates"); v != nil {
		resp.RequestUpdates = v
	}
}

// iconSizeToDevice maps admin/DB enum names to pixel sizes expected by the Android launcher (Java IconSize.getTransmittedValue).
func iconSizeToDevice(value string) (int, bool) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "SMALL":
		return 100, true
	case "MEDIUM":
		return 120, true
	case "LARGE":
		return 140, true
	case "100", "120", "140":
		n, _ := strconv.Atoi(value)
		return n, true
	default:
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil && n > 0 {
			return n, true
		}
	}
	return 0, false
}

func iconSizePtr(m map[string]json.RawMessage, key string) *int {
	if n := intPtr(m, key); n != nil && *n > 0 {
		return n
	}
	if s := stringPtr(m, key); s != nil {
		if n, ok := iconSizeToDevice(*s); ok {
			return &n
		}
	}
	return nil
}

func boolPtr(m map[string]json.RawMessage, key string) *bool {
	raw, ok := m[key]
	if !ok {
		return nil
	}
	var b bool
	if err := json.Unmarshal(raw, &b); err != nil {
		return nil
	}
	return &b
}

func boolPtrTrueOnly(m map[string]json.RawMessage, key string) *bool {
	b := boolPtr(m, key)
	if b != nil && *b {
		return b
	}
	return nil
}

func boolVal(m map[string]json.RawMessage, key string) bool {
	b := boolPtr(m, key)
	return b != nil && *b
}

func stringPtr(m map[string]json.RawMessage, key string) *string {
	raw, ok := m[key]
	if !ok {
		return nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil
	}
	if s == "" {
		return nil
	}
	return &s
}

func intPtr(m map[string]json.RawMessage, key string) *int {
	raw, ok := m[key]
	if !ok {
		return nil
	}
	var n int
	if err := json.Unmarshal(raw, &n); err == nil {
		return &n
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		i := int(f)
		return &i
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if i, err := strconv.Atoi(s); err == nil {
			return &i
		}
	}
	return nil
}
