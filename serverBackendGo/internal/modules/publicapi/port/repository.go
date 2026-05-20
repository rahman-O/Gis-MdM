package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/publicapi/domain"
)

// DeviceRepository unauthenticated device/application access.
type DeviceRepository interface {
	FindDeviceByNumber(ctx context.Context, number string) (*domain.DeviceRef, error)
	CustomerFilesDir(ctx context.Context, customerID int) (string, error)
	HasDuplicateApp(ctx context.Context, customerID int, pkg, version string) (bool, error)
	InsertApplication(ctx context.Context, customerID int, name, pkg, version, url string, flags domain.UploadAppRequest) error
}
