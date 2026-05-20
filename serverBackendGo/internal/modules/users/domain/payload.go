package domain

import authdomain "github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"

// UserPayload is the PUT /private/users body (create/update).
type UserPayload struct {
	ID                  *int64                `json:"id"`
	Login               string                `json:"login"`
	Name                string                `json:"name"`
	Email               string                `json:"email"`
	OldPassword         string                `json:"oldPassword"`
	NewPassword         string                `json:"newPassword"`
	UserRole            *RoleRef              `json:"userRole"`
	AllDevicesAvailable *bool                 `json:"allDevicesAvailable"`
	AllConfigAvailable  *bool                 `json:"allConfigAvailable"`
	Groups              []authdomain.LookupItem `json:"groups"`
	Configurations      []authdomain.LookupItem `json:"configurations"`
}

// ProfilePayload is PUT /details body.
type ProfilePayload struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// RoleRef references a role by id.
type RoleRef struct {
	ID int `json:"id"`
}
