package domain

// UserView is the login response DTO (mirrors Java UserView JSON).
type UserView struct {
	ID                   int64           `json:"id"`
	Login                string          `json:"login"`
	Email                string          `json:"email,omitempty"`
	Name                 string          `json:"name,omitempty"`
	CustomerID           int             `json:"customerId"`
	MasterCustomer       bool            `json:"masterCustomer"`
	Editable             bool            `json:"editable"`
	SingleCustomer       bool            `json:"singleCustomer"`
	UserRole             *UserRoleView   `json:"userRole,omitempty"`
	SuperAdmin           bool            `json:"superAdmin"`
	AllDevicesAvailable  bool            `json:"allDevicesAvailable"`
	AllConfigAvailable   bool            `json:"allConfigAvailable"`
	PasswordReset        bool            `json:"passwordReset"`
	AuthToken            string          `json:"authToken"`
	PasswordResetToken   string          `json:"passwordResetToken,omitempty"`
	Groups               []LookupItem    `json:"groups,omitempty"`
	Configurations       []LookupItem    `json:"configurations,omitempty"`
	TwoFactor            *bool           `json:"twoFactor,omitempty"`
	TwoFactorAccepted    *bool           `json:"twoFactorAccepted,omitempty"`
	IdleLogout           *int            `json:"idleLogout,omitempty"`
}

// UserRoleView wraps role for JSON.
type UserRoleView struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	SuperAdmin  bool              `json:"superAdmin"`
	Permissions []PermissionView  `json:"permissions,omitempty"`
}

// PermissionView exposes permission name.
type PermissionView struct {
	Name string `json:"name"`
}

// NewUserView builds a view from domain user (password cleared).
func NewUserView(u *User) *UserView {
	if u == nil {
		return nil
	}
	v := &UserView{
		ID:                  u.ID,
		Login:               u.Login,
		Email:               u.Email,
		Name:                u.Name,
		CustomerID:          u.CustomerID,
		MasterCustomer:      u.MasterCustomer,
		Editable:            true,
		SingleCustomer:      u.SingleCustomer,
		SuperAdmin:          u.UserRole != nil && u.UserRole.SuperAdmin,
		AllDevicesAvailable: u.AllDevicesAvailable,
		AllConfigAvailable:  u.AllConfigAvailable,
		PasswordReset:       u.PasswordReset,
		AuthToken:           u.AuthToken,
		PasswordResetToken:  u.PasswordResetToken,
		Groups:              u.Groups,
		Configurations:      u.Configurations,
		IdleLogout:          u.IdleLogout,
	}
	if u.TwoFactor {
		t := true
		v.TwoFactor = &t
	}
	if u.TwoFactorAccepted {
		t := true
		v.TwoFactorAccepted = &t
	}
	if u.UserRole != nil {
		perms := make([]PermissionView, 0, len(u.UserRole.Permissions))
		for _, p := range u.UserRole.Permissions {
			if p.Name != "" {
				perms = append(perms, PermissionView{Name: p.Name})
			}
		}
		v.UserRole = &UserRoleView{
			ID:          u.UserRole.ID,
			Name:        u.UserRole.Name,
			SuperAdmin:  u.UserRole.SuperAdmin,
			Permissions: perms,
		}
	}
	return v
}
