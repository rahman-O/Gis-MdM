package domain

// PushMessageFilter for plugin search.
type PushMessageFilter struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

// PluginPushMessage history row.
type PluginPushMessage struct {
	ID          int64  `json:"id"`
	CustomerID  int64  `json:"customerId"`
	DeviceID    int64  `json:"deviceId"`
	Ts          int64  `json:"ts"`
	MessageType string `json:"messageType"`
	Payload     string `json:"payload"`
}

// PushSendRequest for plugin send.
type PushSendRequest struct {
	Scope        string `json:"scope"`
	DeviceNumber string `json:"deviceNumber"`
	GroupID      int64  `json:"groupId"`
	MessageType  string `json:"messageType"`
	Payload      string `json:"payload"`
}

// PaginatedMessages wraps search results.
type PaginatedMessages struct {
	Items []PluginPushMessage `json:"items"`
	Total int64               `json:"totalItems"`
}

type PushScheduleFilter struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

type PluginPushSchedule struct {
	ID              int64  `json:"id"`
	CustomerID      int64  `json:"customerId"`
	DeviceID        int64  `json:"deviceId"`
	GroupID         int64  `json:"groupId"`
	ConfigurationID int64  `json:"configurationId"`
	Scope           string `json:"scope"`
	DeviceNumber    string `json:"deviceNumber"`
	MessageType     string `json:"messageType"`
	Payload         string `json:"payload"`
	Comment         string `json:"comment"`
	Min             string `json:"min"`
	Hour            string `json:"hour"`
	Day             string `json:"day"`
	Weekday         string `json:"weekday"`
	Month           string `json:"month"`
}

type PaginatedSchedules struct {
	Items []PluginPushSchedule `json:"items"`
	Total int64                `json:"totalItemsCount"`
}
