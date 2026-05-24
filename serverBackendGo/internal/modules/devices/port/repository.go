package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/devices/domain"
)

// UserScope is loaded per request for device visibility.
type UserScope struct {
	UserID              int64
	CustomerID          int
	AllDevicesAvailable bool
}

// DeviceRepository persists devices.
type DeviceRepository interface {
	LoadUserScope(ctx context.Context, userID int64) (*UserScope, error)
	Search(ctx context.Context, scope UserScope, req domain.SearchRequest) ([]domain.DeviceView, error)
	Count(ctx context.Context, scope UserScope, req domain.SearchRequest) (int64, error)
	ListConfigurations(ctx context.Context, customerID int) (map[int]domain.ConfigurationView, error)
	GetByNumber(ctx context.Context, scope UserScope, number string) (*domain.DeviceView, error)
	GetByID(ctx context.Context, customerID int, id int) (*domain.DeviceView, error)
	ExistsNumber(ctx context.Context, customerID int, number string, excludeID int) (bool, error)
	CountDevices(ctx context.Context, customerID int) (int64, error)
	DeviceLimit(ctx context.Context, customerID int) (int, error)
	Insert(ctx context.Context, customerID int, d domain.SaveDevice) (int, error)
	Update(ctx context.Context, customerID int, d domain.SaveDevice) error
	UpdateConfigurationBulk(ctx context.Context, customerID int, ids []int, configID int) error
	Delete(ctx context.Context, customerID int, id int) error
	DeleteBulk(ctx context.Context, customerID int, ids []int) error
	UpdateGroupBulk(ctx context.Context, customerID int, req domain.GroupBulkRequest) error
	Autocomplete(ctx context.Context, scope UserScope, filter string, limit int) ([]domain.LookupItem, error)
	UpdateDescription(ctx context.Context, customerID int, id int, description string) error
	ListAppSettings(ctx context.Context, deviceID int) ([]domain.AppSetting, error)
	SaveAppSettings(ctx context.Context, deviceID int, settings []domain.AppSetting) error
	MoveTreeNode(ctx context.Context, customerID int, deviceID int, treeNodeID int) error
}
