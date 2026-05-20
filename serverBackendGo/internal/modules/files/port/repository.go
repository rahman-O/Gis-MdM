package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/applications/domain"
	filesdomain "github.com/gis-mdm/server-backend-go/internal/modules/files/domain"
)

// FileRepository persists uploaded file metadata.
type FileRepository interface {
	List(ctx context.Context, customerID int, filter string) ([]filesdomain.UploadedFile, error)
	GetByID(ctx context.Context, customerID, id int) (*filesdomain.UploadedFile, error)
	Insert(ctx context.Context, f *filesdomain.UploadedFile) error
	Update(ctx context.Context, f *filesdomain.UploadedFile) error
	Delete(ctx context.Context, id int) error
	IsUsedByConfiguration(ctx context.Context, fileID int) (bool, error)
	IsUsedByIcon(ctx context.Context, fileID int) (bool, error)
	UsingConfigurationNames(ctx context.Context, customerID, fileID int) ([]string, error)
	UsingIconNames(ctx context.Context, customerID, fileID int) ([]string, error)
	GetFileConfigurations(ctx context.Context, customerID, userID, fileID int) ([]filesdomain.FileConfigurationLink, error)
	DeleteConfigurationFile(ctx context.Context, linkID int) error
	InsertConfigurationFile(ctx context.Context, configurationID, fileID int, devicePath string) error
	CountByPath(ctx context.Context, customerID int, id *int, filePath string) (int64, error)
}

// CustomerRepository loads per-tenant file settings.
type CustomerRepository interface {
	GetMeta(ctx context.Context, customerID int) (*filesdomain.CustomerMeta, error)
	CountCustomers(ctx context.Context) (int, error)
}

// ApplicationLookup finds apps referencing a file URL.
type ApplicationLookup interface {
	SearchByURL(ctx context.Context, customerID int, url string) ([]domain.Application, error)
	FindVersionByPkgCode(ctx context.Context, customerID int, pkg string, versionCode int) (*domain.ApplicationVersion, error)
	FindVersionByPkgVersion(ctx context.Context, customerID int, pkg, version string) (*domain.ApplicationVersion, error)
	FindAppsByPkg(ctx context.Context, customerID int, pkg string) ([]domain.Application, error)
}

// PushNotifier stub for configuration notify.
type PushNotifier interface {
	NotifyConfigurationUpdate(configurationID int)
}

type noopPush struct{}

func (noopPush) NotifyConfigurationUpdate(int) {}

// NoopPush returns a no-op notifier.
func NoopPush() PushNotifier { return noopPush{} }
