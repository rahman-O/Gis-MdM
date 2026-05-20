package domain

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

// DeviceView mirrors React DeviceView (subset for Phase 4).
type DeviceView struct {
	ID              int          `json:"id"`
	Number          string       `json:"number"`
	Description     *string      `json:"description"`
	LastUpdate      *int64       `json:"lastUpdate"`
	ConfigurationID *int         `json:"configurationId"`
	IMEI            *string      `json:"imei"`
	Phone           *string      `json:"phone"`
	StatusCode      *string      `json:"statusCode"`
	Groups          []LookupItem `json:"groups"`
	Custom1         *string      `json:"custom1"`
	Custom2         *string      `json:"custom2"`
	Custom3         *string      `json:"custom3"`
	OldNumber       *string      `json:"oldNumber"`
}

// SearchRequest mirrors DeviceSearchRequest (React uses pageNum).
type SearchRequest struct {
	PageNum         int     `json:"pageNum"`
	PageSize        int     `json:"pageSize"`
	Value           *string `json:"value"`
	GroupID         *int    `json:"groupId"`
	ConfigurationID *int    `json:"configurationId"`
	SortBy          *string `json:"sortBy"`
	SortDir         *string `json:"sortDir"`
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
