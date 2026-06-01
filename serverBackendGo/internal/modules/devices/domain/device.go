package domain

import "strings"

// LookupItem mirrors com.hmdm.rest.json.LookupItem.
type LookupItem struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
}

// ConfigurationView is a minimal configuration entry for device search.
type ConfigurationView struct {
	ID             int     `json:"id"`
	Name           *string `json:"name"`
	PermissiveMode *bool   `json:"permissiveMode"`
}

// DevicePermissions mirrors React DevicePermissions (subset).
type DevicePermissions struct {
	PermissionStatus *string `json:"permissionStatus"`
	Details          *string `json:"details"`
}

// DeviceApplication mirrors installed app telemetry in device info.
type DeviceApplication struct {
	Pkg     *string `json:"pkg"`
	Version *string `json:"version"`
	Status  *string `json:"status"`
}

// DeviceFile mirrors configuration file install status in device info.
type DeviceFile struct {
	Path   *string `json:"path"`
	Status *string `json:"status"`
}

// DeviceInfoView mirrors React DeviceInfoView (agent telemetry).
type DeviceInfoView struct {
	BatteryLevel    *int                 `json:"batteryLevel"`
	Model           *string              `json:"model"`
	AndroidVersion  *string              `json:"androidVersion"`
	Serial          *string              `json:"serial"`
	IMEI            *string              `json:"imei"`
	Phone           *string              `json:"phone"`
	Location        *string              `json:"location"`
	Permissions     *DevicePermissions   `json:"permissions"`
	Applications    []DeviceApplication  `json:"applications"`
	Files           []DeviceFile         `json:"files"`
	DefaultLauncher *bool                `json:"defaultLauncher"`
	MdmMode         *bool                `json:"mdmMode"`
	KioskMode       *bool                `json:"kioskMode"`
	EnrollTime      *int64               `json:"enrollTime"`
	PublicIP        *string              `json:"publicIp"`
	LauncherVersion *string              `json:"launcherVersion"`
}

// DeviceView mirrors React DeviceView (subset for Phase 4 + 012 telemetry).
type DeviceView struct {
	ID              int          `json:"id"`
	Number          string       `json:"number"`
	Description     *string      `json:"description"`
	LastUpdate      *int64       `json:"lastUpdate"`
	ConfigurationID *int         `json:"configurationId"`
	TreeNodeID      *int         `json:"treeNodeId,omitempty"`
	AgentID         *string      `json:"agentId,omitempty"`
	EnrollmentRouteID *int       `json:"enrollmentRouteId,omitempty"`
	EnrollmentState *string      `json:"enrollmentState,omitempty"`
	IMEI            *string      `json:"imei"`
	Phone           *string      `json:"phone"`
	Model           *string      `json:"model"`
	BatteryLevel    *int         `json:"batteryLevel"`
	AndroidVersion  *string      `json:"androidVersion"`
	Serial          *string      `json:"serial"`
	LauncherVersion *string      `json:"launcherVersion"`
	StatusCode      *string      `json:"statusCode"`
	Groups          []LookupItem `json:"groups"`
	Custom1         *string      `json:"custom1"`
	Custom2         *string      `json:"custom2"`
	Custom3         *string      `json:"custom3"`
	OldNumber       *string      `json:"oldNumber"`
	Info            *DeviceInfoView `json:"info,omitempty"`
}

// SearchRequest mirrors DeviceSearchRequest (React uses pageNum).
type SearchRequest struct {
	PageNum         int     `json:"pageNum"`
	PageSize        int     `json:"pageSize"`
	Value           *string `json:"value"`
	GroupID            *int    `json:"groupId"`
	ConfigurationID    *int    `json:"configurationId"`
	TreeNodeID         *int    `json:"treeNodeId"`
	IncludeDescendants *bool   `json:"includeDescendants"`
	Status             *string `json:"status"`
	AndroidVersion  *string `json:"androidVersion"`
	LauncherVersion *string `json:"launcherVersion"`
	MdmMode         *bool   `json:"mdmMode"`
	KioskMode       *bool   `json:"kioskMode"`
	InstallationStatus *string `json:"installationStatus"`
	SortBy          *string `json:"sortBy"`
	SortDir         *string `json:"sortDir"`
	DateFrom        *int64  `json:"dateFrom"`
	DateTo          *int64  `json:"dateTo"`
	OnlineEarlierMillis *int64 `json:"onlineEarlierMillis"`
	OnlineLaterMillis   *int64 `json:"onlineLaterMillis"`
	EnrollmentDateFrom  *int64 `json:"enrollmentDateFrom"`
	EnrollmentDateTo    *int64 `json:"enrollmentDateTo"`
	ImeiChanged     *bool   `json:"imeiChanged"`
	FastSearch      *bool   `json:"fastSearch"`
}

// Normalize applies defaults for paging.
func (r *SearchRequest) Normalize() {
	if r.PageNum < 1 {
		r.PageNum = 1
	}
	if r.PageSize < 1 {
		r.PageSize = 50
	}
	if r.PageSize > 500 {
		r.PageSize = 500
	}
}

// Prepare normalizes paging and search value wildcards (Java DeviceSearchRequest.getValue).
func (r *SearchRequest) Prepare() {
	r.Normalize()
	if r.Value == nil {
		return
	}
	v := strings.TrimSpace(*r.Value)
	if v == "" {
		r.Value = nil
		return
	}
	fast := r.FastSearch != nil && *r.FastSearch
	if !fast {
		wrapped := "%" + v + "%"
		r.Value = &wrapped
	} else {
		r.Value = &v
	}
}

// DevicePage is the paginated devices block inside DeviceListView.
type DevicePage struct {
	Items           []DeviceView `json:"items"`
	TotalItemsCount int64        `json:"totalItemsCount"`
}

// DeviceListView mirrors DeviceListResponse for React.
type DeviceListView struct {
	Configurations map[int]ConfigurationView `json:"configurations"`
	Devices        DevicePage                `json:"devices"`
}

// SaveDevice is PUT /devices body.
type SaveDevice struct {
	ID              *int         `json:"id"`
	IDs             []int        `json:"ids"`
	Number          *string      `json:"number"`
	Description     *string      `json:"description"`
	ConfigurationID *int         `json:"configurationId"`
	Groups          []LookupItem `json:"groups"`
}

// BulkDeleteRequest is POST /deleteBulk.
type BulkDeleteRequest struct {
	IDs []int `json:"ids"`
}

// MoveTreeRequest is POST /devices/:id/move-tree.
type MoveTreeRequest struct {
	TreeNodeID int `json:"treeNodeId"`
}

// GroupBulkRequest is POST /groupBulk.
type GroupBulkRequest struct {
	IDs    []int        `json:"ids"`
	Action string       `json:"action"`
	Groups []LookupItem `json:"groups"`
}

// AppSetting mirrors per-device application settings.
type AppSetting struct {
	ApplicationPkg *string `json:"applicationPkg"`
	Name           *string `json:"name"`
	Type           *string `json:"type"`
	Value          *string `json:"value"`
}
