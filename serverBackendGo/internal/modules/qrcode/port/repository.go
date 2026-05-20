package port

import "context"

// ConfigByKey resolves configuration and main app version for QR.
type ConfigByKey interface {
	ConfigurationByQRKey(ctx context.Context, key string) (*QRConfig, error)
}

type QRConfig struct {
	ID           int64
	Name         string
	LauncherURL  string
	MainAppPkg   string
	MainAppURL   string
	AdminExtras  string
	CustomerID   int64
	FilesDir     string
}
