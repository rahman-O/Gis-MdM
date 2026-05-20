package domain

// User is the authenticated account (persistence aggregate).
type User struct {
	ID                   int64
	Login                string
	Email                string
	Name                 string
	Password             string
	CustomerID           int
	MasterCustomer       bool
	AllDevicesAvailable  bool
	AllConfigAvailable   bool
	PasswordReset        bool
	AuthToken            string
	PasswordResetToken   string
	TwoFactor            bool
	TwoFactorAccepted    bool
	LastLoginFail        int64
	IdleLogout           *int
	SingleCustomer       bool
	UserRole             *UserRole
	Groups               []LookupItem
	Configurations       []LookupItem
}

// UserRole holds role and permissions.
type UserRole struct {
	ID          int
	Name        string
	SuperAdmin  bool
	Permissions []Permission
}

// Permission is a named capability.
type Permission struct {
	ID         int
	Name       string
	SuperAdmin bool
}

// LookupItem is id/name pair for groups and configurations.
type LookupItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CustomerSettings subset for login.
type CustomerSettings struct {
	TwoFactor   bool
	IdleLogout  *int
}
