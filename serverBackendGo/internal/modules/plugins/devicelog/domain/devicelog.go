package domain

type Settings struct {
	ID                int64 `json:"id"`
	CustomerID        int64 `json:"customerId"`
	LogsPreservePeriod int  `json:"logsPreservePeriod"`
}

type Rule struct {
	ID              int64  `json:"id"`
	SettingID       int64  `json:"settingId"`
	Name            string `json:"name"`
	Active          bool   `json:"active"`
	ApplicationID   int64  `json:"applicationId"`
	Severity        string `json:"severity"`
	Filter          string `json:"filter,omitempty"`
	GroupID         *int64 `json:"groupId,omitempty"`
	ConfigurationID *int64 `json:"configurationId,omitempty"`
}

type LogFilter struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	DeviceID int64  `json:"deviceId"`
	Severity string `json:"severity"`
}

type LogRecord struct {
	ID            int64  `json:"id"`
	CreateTime    int64  `json:"createTime"`
	DeviceID      int64  `json:"deviceId"`
	ApplicationID int64  `json:"applicationId"`
	Severity      string `json:"severity"`
	Message       string `json:"message"`
}

type UploadRecord struct {
	Severity  string `json:"severity"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type PaginatedLogs struct {
	Items []LogRecord `json:"items"`
	Total int64       `json:"totalItemsCount"`
}
