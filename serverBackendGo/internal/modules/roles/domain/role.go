package domain

// Permission is a named capability.
type Permission struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	SuperAdmin  bool   `json:"superAdmin,omitempty"`
}

// Role is a user role with permissions.
type Role struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	SuperAdmin  bool         `json:"superAdmin"`
	Permissions []Permission `json:"permissions,omitempty"`
}

// RolePayload is PUT /roles body.
type RolePayload struct {
	ID          *int         `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Permissions []Permission `json:"permissions"`
}
