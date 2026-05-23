package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
)

// ConfigRepository persists configurations and related rows.
type ConfigRepository interface {
	ListByCustomer(ctx context.Context, customerID int) ([]domain.LookupItem, error)
	Search(ctx context.Context, customerID int) ([]domain.Configuration, error)
	SearchByValue(ctx context.Context, customerID int, value string) ([]domain.Configuration, error)
	GetByID(ctx context.Context, customerID, id int) (*domain.Configuration, error)
	GetByName(ctx context.Context, customerID int, name string) (*domain.Configuration, error)
	CountDevicesUsing(ctx context.Context, configurationID int) (int64, error)
	Insert(ctx context.Context, customerID int, cfg domain.Configuration) (int, error)
	Update(ctx context.Context, customerID int, cfg domain.Configuration) error
	Delete(ctx context.Context, customerID, id int) error
	Copy(ctx context.Context, customerID int, req domain.CopyRequest) (int, error)
	ListAllApplicationsForPicker(ctx context.Context, customerID int) ([]domain.ConfigurationApplication, error)
	ListConfigurationApplications(ctx context.Context, customerID, configurationID int) ([]domain.ConfigurationApplication, error)
	UpgradeApplication(ctx context.Context, customerID int, req domain.UpgradeApplicationRequest) error
}
