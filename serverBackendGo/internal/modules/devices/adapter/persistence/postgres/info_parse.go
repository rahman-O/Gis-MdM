package postgres

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/devices/domain"
)

type deviceInfoJSON struct {
	BatteryLevel    *int                    `json:"batteryLevel"`
	Model           *string                 `json:"model"`
	AndroidVersion  *string                 `json:"androidVersion"`
	Serial          *string                 `json:"serial"`
	IMEI            *string                 `json:"imei"`
	Phone           *string                 `json:"phone"`
	DefaultLauncher *bool                   `json:"defaultLauncher"`
	MdmMode         *bool                   `json:"mdmMode"`
	KioskMode       *bool                   `json:"kioskMode"`
	LauncherVersion *string                 `json:"launcherVersion"`
	Applications    []domain.DeviceApplication `json:"applications"`
	Files           []domain.DeviceFile        `json:"files"`
	Permissions     []int                   `json:"permissions"`
	Location        json.RawMessage         `json:"location"`
}

func parseDeviceInfo(info sql.NullString, infojson []byte, enrollTime sql.NullInt64, publicIP sql.NullString) *domain.DeviceInfoView {
	var merged deviceInfoJSON
	if info.Valid && strings.TrimSpace(info.String) != "" {
		_ = json.Unmarshal([]byte(info.String), &merged)
	}
	if len(infojson) > 0 {
		var fromJSON deviceInfoJSON
		if err := json.Unmarshal(infojson, &fromJSON); err == nil {
			mergeDeviceInfo(&merged, &fromJSON)
		}
	}
	if merged.BatteryLevel == nil && merged.Model == nil && merged.AndroidVersion == nil &&
		len(merged.Applications) == 0 && len(merged.Files) == 0 {
		return nil
	}
	out := &domain.DeviceInfoView{
		BatteryLevel:    merged.BatteryLevel,
		Model:           merged.Model,
		AndroidVersion:  merged.AndroidVersion,
		Serial:          merged.Serial,
		IMEI:            merged.IMEI,
		Phone:           merged.Phone,
		DefaultLauncher: merged.DefaultLauncher,
		MdmMode:         merged.MdmMode,
		KioskMode:       merged.KioskMode,
		LauncherVersion: merged.LauncherVersion,
	}
	if len(merged.Applications) > 0 {
		out.Applications = merged.Applications
	}
	if len(merged.Files) > 0 {
		out.Files = merged.Files
	}
	if len(merged.Permissions) > 0 {
		out.Permissions = permissionsFromInts(merged.Permissions)
	}
	if len(merged.Location) > 0 {
		loc := string(merged.Location)
		if loc != "null" {
			out.Location = &loc
		}
	}
	if enrollTime.Valid {
		out.EnrollTime = &enrollTime.Int64
	}
	if publicIP.Valid {
		out.PublicIP = &publicIP.String
	}
	return out
}

func mergeDeviceInfo(dst, src *deviceInfoJSON) {
	if dst == nil || src == nil {
		return
	}
	if src.BatteryLevel != nil {
		dst.BatteryLevel = src.BatteryLevel
	}
	if src.Model != nil {
		dst.Model = src.Model
	}
	if src.AndroidVersion != nil {
		dst.AndroidVersion = src.AndroidVersion
	}
	if src.Serial != nil {
		dst.Serial = src.Serial
	}
	if src.IMEI != nil {
		dst.IMEI = src.IMEI
	}
	if src.Phone != nil {
		dst.Phone = src.Phone
	}
	if src.DefaultLauncher != nil {
		dst.DefaultLauncher = src.DefaultLauncher
	}
	if src.MdmMode != nil {
		dst.MdmMode = src.MdmMode
	}
	if src.KioskMode != nil {
		dst.KioskMode = src.KioskMode
	}
	if src.LauncherVersion != nil {
		dst.LauncherVersion = src.LauncherVersion
	}
	if len(src.Applications) > 0 {
		dst.Applications = src.Applications
	}
	if len(src.Files) > 0 {
		dst.Files = src.Files
	}
	if len(src.Permissions) > 0 {
		dst.Permissions = src.Permissions
	}
	if len(src.Location) > 0 {
		dst.Location = src.Location
	}
}

func permissionsFromInts(perms []int) *domain.DevicePermissions {
	if len(perms) == 0 {
		return nil
	}
	parts := make([]string, len(perms))
	for i, p := range perms {
		parts[i] = strconv.Itoa(p)
	}
	details := strings.Join(parts, ",")
	status := "unknown"
	return &domain.DevicePermissions{
		PermissionStatus: &status,
		Details:          &details,
	}
}
