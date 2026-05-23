package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/applications/domain"
)

// ApplicationRepository persists applications and links.
type ApplicationRepository interface {
	Search(ctx context.Context, customerID int) ([]domain.Application, error)
	SearchByValue(ctx context.Context, customerID int, value string) ([]domain.Application, error)
	GetByID(ctx context.Context, customerID, id int) (*domain.Application, error)
	ListVersions(ctx context.Context, customerID, applicationID int) ([]domain.ApplicationVersion, error)
	SaveAndroid(ctx context.Context, customerID int, app domain.Application) (*domain.Application, error)
	SaveWeb(ctx context.Context, customerID int, app domain.Application) (*domain.Application, error)
	SaveVersion(ctx context.Context, customerID int, ver domain.ApplicationVersion) (*domain.ApplicationVersion, error)
	DeleteApp(ctx context.Context, customerID, id int) error
	DeleteVersion(ctx context.Context, customerID, versionID int) error
	ValidatePkg(ctx context.Context, customerID int, req domain.ValidatePkgRequest) ([]domain.Application, error)
	GetAppConfigurations(ctx context.Context, customerID, applicationID int) ([]domain.ApplicationConfigurationLink, error)
	UpdateAppConfigurations(ctx context.Context, customerID int, req domain.LinkConfigurationsToAppRequest) error
	GetVersionConfigurations(ctx context.Context, customerID, versionID int) ([]domain.ApplicationVersionConfigurationLink, error)
	UpdateVersionConfigurations(ctx context.Context, customerID int, req domain.LinkConfigurationsToAppVersionRequest) error
	AdminSearch(ctx context.Context, value string) ([]domain.Application, error)
	TurnIntoCommon(ctx context.Context, id int) error
	CustomerFilesDir(ctx context.Context, customerID int) (string, error)
}
