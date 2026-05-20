package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/sync/domain"
)

// SyncRepository loads devices and builds configuration payloads.
type SyncRepository interface {
	FindByNumber(ctx context.Context, number string) (*domain.DeviceRecord, error)
	FindByOldNumber(ctx context.Context, number string) (*domain.DeviceRecord, error)
	CreateOnDemand(ctx context.Context, number string, opts domain.DeviceCreateOptions, defaultCustomerID int64) (*domain.DeviceRecord, error)
	CompleteMigration(ctx context.Context, deviceID int64) error
	TouchLastUpdate(ctx context.Context, deviceID int64) error
	UpdateInfo(ctx context.Context, deviceID int64, infoJSON string, publicIP string) error
	UpdateCustomProps(ctx context.Context, deviceID int64, custom1, custom2, custom3 *string) error
	SaveApplicationSettings(ctx context.Context, deviceID int64, settings []domain.SyncApplicationSetting) error
	BuildSyncResponse(ctx context.Context, dev domain.DeviceRecord, baseURL, filesDir, cpuArch, mobileName, vendor string) (*domain.SyncResponse, error)
	CountCustomers(ctx context.Context) (int, error)
}
