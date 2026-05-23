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
	NewDeviceGroupID         *int    `json:"newDeviceGroupId,omitempty"`
	PhoneNumberFormat        string  `json:"phoneNumberFormat,omitempty"`
	CustomPropertyName1      string  `json:"customPropertyName1,omitempty"`
	CustomPropertyName2      string  `json:"customPropertyName2,omitempty"`
	CustomPropertyName3      string  `json:"customPropertyName3,omitempty"`
	CustomMultiline1         bool    `json:"customMultiline1"`
	CustomMultiline2         bool    `json:"customMultiline2"`
	CustomMultiline3         bool    `json:"customMultiline3"`
	CustomSend1              bool    `json:"customSend1"`
	CustomSend2              bool    `json:"customSend2"`
	CustomSend3              bool    `json:"customSend3"`
	DesktopHeaderTemplate    string  `json:"desktopHeaderTemplate,omitempty"`
	SendDescription          bool    `json:"sendDescription"`
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
	RoleID                                int  `json:"roleId"`
	CustomerID                            int  `json:"customerId,omitempty"`
	ColumnDisplayedDeviceStatus           bool `json:"columnDisplayedDeviceStatus"`
	ColumnDisplayedDeviceDate             bool `json:"columnDisplayedDeviceDate"`
	ColumnDisplayedDeviceNumber           bool `json:"columnDisplayedDeviceNumber"`
	ColumnDisplayedDeviceModel            bool `json:"columnDisplayedDeviceModel"`
	ColumnDisplayedDevicePermissionsStatus bool `json:"columnDisplayedDevicePermissionsStatus"`
	ColumnDisplayedDeviceAppInstallStatus bool `json:"columnDisplayedDeviceAppInstallStatus"`
	ColumnDisplayedDeviceConfiguration    bool `json:"columnDisplayedDeviceConfiguration"`
	ColumnDisplayedDeviceImei             bool `json:"columnDisplayedDeviceImei"`
	ColumnDisplayedDevicePhone            bool `json:"columnDisplayedDevicePhone"`
	ColumnDisplayedDeviceDesc             bool `json:"columnDisplayedDeviceDesc"`
	ColumnDisplayedDeviceGroup            bool `json:"columnDisplayedDeviceGroup"`
	ColumnDisplayedLauncherVersion        bool `json:"columnDisplayedLauncherVersion"`
	ColumnDisplayedDeviceFilesStatus      bool `json:"columnDisplayedDeviceFilesStatus"`
	ColumnDisplayedBatteryLevel           bool `json:"columnDisplayedBatteryLevel"`
	ColumnDisplayedDefaultLauncher        bool `json:"columnDisplayedDefaultLauncher"`
	ColumnDisplayedCustom1                bool `json:"columnDisplayedCustom1"`
	ColumnDisplayedCustom2                bool `json:"columnDisplayedCustom2"`
	ColumnDisplayedCustom3                bool `json:"columnDisplayedCustom3"`
	ColumnDisplayedMdmMode                bool `json:"columnDisplayedMdmMode"`
	ColumnDisplayedKioskMode              bool `json:"columnDisplayedKioskMode"`
	ColumnDisplayedAndroidVersion         bool `json:"columnDisplayedAndroidVersion"`
	ColumnDisplayedEnrollmentDate         bool `json:"columnDisplayedEnrollmentDate"`
	ColumnDisplayedSerial                 bool `json:"columnDisplayedSerial"`
	ColumnDisplayedPublicIp               bool `json:"columnDisplayedPublicIp"`
}

// DefaultUserRoleSettings returns Java-like defaults (all columns visible).
func DefaultUserRoleSettings(roleID, customerID int) UserRoleSettings {
	return UserRoleSettings{
		RoleID: roleID, CustomerID: customerID,
		ColumnDisplayedDeviceStatus: true, ColumnDisplayedDeviceDate: true,
		ColumnDisplayedDeviceNumber: true, ColumnDisplayedDeviceModel: true,
		ColumnDisplayedDevicePermissionsStatus: true, ColumnDisplayedDeviceAppInstallStatus: true,
		ColumnDisplayedDeviceConfiguration: true, ColumnDisplayedDeviceImei: true,
		ColumnDisplayedDevicePhone: true, ColumnDisplayedDeviceDesc: true,
		ColumnDisplayedDeviceGroup: true, ColumnDisplayedLauncherVersion: true,
		ColumnDisplayedDeviceFilesStatus: true, ColumnDisplayedBatteryLevel: true,
		ColumnDisplayedDefaultLauncher: true, ColumnDisplayedCustom1: true,
		ColumnDisplayedCustom2: true, ColumnDisplayedCustom3: true,
		ColumnDisplayedMdmMode: true, ColumnDisplayedKioskMode: true,
		ColumnDisplayedAndroidVersion: true, ColumnDisplayedEnrollmentDate: true,
		ColumnDisplayedSerial: true, ColumnDisplayedPublicIp: true,
	}
}
