package domain

const (
	StatusPending   = 0
	StatusDelivered = 1
	StatusRead      = 2
)

type MessageFilter struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

type Message struct {
	ID         int64  `json:"id"`
	CustomerID int64  `json:"customerId"`
	DeviceID   int64  `json:"deviceId"`
	Ts         int64  `json:"ts"`
	Message    string `json:"message"`
	Status     int    `json:"status"`
}

type SendRequest struct {
	Scope             string `json:"scope"`
	DeviceNumber      string `json:"deviceNumber"`
	GroupID           int64  `json:"groupId"`
	ConfigurationID   int64  `json:"configurationId"`
	Message           string `json:"message"`
}

type PaginatedMessages struct {
	Items []Message `json:"items"`
	Total int64     `json:"totalItemsCount"`
}
