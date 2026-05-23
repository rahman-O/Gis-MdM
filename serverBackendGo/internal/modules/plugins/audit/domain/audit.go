package domain

type AuditLogFilter struct {
	PageNum       int    `json:"pageNum"`
	PageSize      int    `json:"pageSize"`
	DateFrom      *int64 `json:"dateFrom"`
	DateTo        *int64 `json:"dateTo"`
	UserFilter    string `json:"userFilter"`
	MessageFilter string `json:"messageFilter"`
}

type AuditLogRecord struct {
	ID         int64  `json:"id"`
	CreateTime int64  `json:"createTime"`
	CustomerID int64  `json:"customerId"`
	UserID     *int64 `json:"userId,omitempty"`
	Login      string `json:"login,omitempty"`
	Action     string `json:"action,omitempty"`
	Payload    string `json:"payload,omitempty"`
	IPAddress  string `json:"ipAddress,omitempty"`
	ErrorCode  int    `json:"errorCode"`
}

type PaginatedAudit struct {
	Items []AuditLogRecord `json:"items"`
	Total int64            `json:"totalItemsCount"`
}
