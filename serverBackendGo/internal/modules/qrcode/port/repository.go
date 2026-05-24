package port

import "context"

// ConfigByKey resolves configuration and main app version for QR.
type ConfigByKey interface {
	ConfigurationByQRKey(ctx context.Context, key string) (*QRConfig, error)
	CountCustomers(ctx context.Context) (int, error)
}

type QRConfig struct {
	ID                      int64
	Name                    string
	QRCodeKey               string
	CustomerID              int64
	CustomerName            string
	FilesDir                string
	MainAppVersionID        int64
	MainAppPkg              string
	MainAppURL              string
	MainAppFilePath         string
	AppLevelURL             string
	ApkHash                 string
	LauncherURL             string
	EventReceivingComponent string
	AdminExtras             string
	QRParameters            string
	WifiSSID                string
	WifiPassword            string
	WifiSecurityType        string
	MobileEnrollment        bool
	EncryptDevice           bool
	DefaultDeviceIDMode     string
}
