package domain

type Settings struct {
	ID                 int64 `json:"id"`
	CustomerID         int64 `json:"customerId"`
	DataPreservePeriod int   `json:"dataPreservePeriod"`
	SendData           bool  `json:"sendData"`
	IntervalMins       int   `json:"intervalMins"`
}

type DynamicInfo struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
	Ts        int64  `json:"ts,omitempty"`
}

type DeviceDetail struct {
	DeviceNumber string        `json:"deviceNumber"`
	Records      []ParamsRecord `json:"records,omitempty"`
}

type ParamsRecord struct {
	ID           int64  `json:"id"`
	Ts           int64  `json:"ts"`
	BatteryLevel *int   `json:"batteryLevel,omitempty"`
}

type DynamicSearchFilter struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	DeviceID int64  `json:"deviceId"`
}

type PaginatedDynamic struct {
	Items []ParamsRecord `json:"items"`
	Total int64          `json:"totalItemsCount"`
}
