package domain

import "errors"

// Customer is a tenant account (JSON mirrors Java Customer).
type Customer struct {
	ID                    *int    `json:"id,omitempty"`
	Name                  string  `json:"name"`
	Email                 string  `json:"email,omitempty"`
	Description           string  `json:"description,omitempty"`
	FilesDir              string  `json:"filesDir,omitempty"`
	Master                bool    `json:"master,omitempty"`
	Prefix                string  `json:"prefix,omitempty"`
	LastLoginTime         *int64  `json:"lastLoginTime,omitempty"`
	RegistrationTime      *int64  `json:"registrationTime,omitempty"`
	AccountType           *int    `json:"accountType,omitempty"`
	CustomerStatus        string  `json:"customerStatus,omitempty"`
	ExpiryTime            *int64  `json:"expiryTime,omitempty"`
	DeviceLimit           *int    `json:"deviceLimit,omitempty"`
	DeviceConfigurationID *int    `json:"deviceConfigurationId,omitempty"`
	MainUserName          string  `json:"mainUserName,omitempty"`
}

// SearchRequest mirrors CustomerSearchRequest.
type SearchRequest struct {
	CurrentPage    int     `json:"currentPage"`
	PageSize       int     `json:"pageSize"`
	SearchValue    *string `json:"searchValue"`
	SortValue      *string `json:"sortValue"`
	SortDirection  *string `json:"sortDirection"`
	AccountType    *int    `json:"accountType"`
	CustomerStatus *string `json:"customerStatus"`
}

// Paginated is the API page wrapper (Java PaginatedData).
type Paginated struct {
	Items            []Customer `json:"items"`
	TotalItemsCount  int64      `json:"totalItemsCount"`
}

// NormalizeSearch applies defaults for paging.
func (r *SearchRequest) NormalizeSearch() {
	if r.CurrentPage < 1 {
		r.CurrentPage = 1
	}
	if r.PageSize < 1 {
		r.PageSize = 100
	}
}

var ErrNotSuperAdmin = errors.New("error.permission.denied")
