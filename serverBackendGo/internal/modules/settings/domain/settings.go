package domain

// Settings is the tenant settings DTO (subset of Java Settings JSON).
type Settings struct {
	ID                       int     `json:"id"`
	CustomerID               int     `json:"customerId"`
	CustomerName             string  `json:"customerName,omitempty"`
	SingleCustomer           bool    `json:"singleCustomer"`
	Language                 string  `json:"language,omitempty"`
	UseDefaultLanguage       bool    `json:"useDefaultLanguage"`
	CreateNewDevices         bool    `json:"createNewDevices"`
	NewDeviceConfigurationID *int    `json:"newDeviceConfigurationId,omitempty"`
	PasswordLength           int     `json:"passwordLength"`
	PasswordStrength         int     `json:"passwordStrength"`
	TwoFactor                bool    `json:"twoFactor"`
	IdleLogout               *int    `json:"idleLogout,omitempty"`
	BackgroundColor          string  `json:"backgroundColor,omitempty"`
	TextColor                string  `json:"textColor,omitempty"`
	BackgroundImageURL       string  `json:"backgroundImageUrl,omitempty"`
	IconSize                 string  `json:"iconSize,omitempty"`
	DesktopHeader            string  `json:"desktopHeader,omitempty"`
	UnsecureEnrollment       bool    `json:"unsecureEnrollment"`
	DeviceFastSearch         bool    `json:"deviceFastSearch"`
	SendDeviceInfoExpiryDays int     `json:"sendDeviceInfoExpiryDays"`
}

// UserRoleSettings mirrors Java UserRoleSettings for column prefs.
type UserRoleSettings struct {
	RoleID int `json:"roleId"`
}
