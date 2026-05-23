package application

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/applications/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type appRepoStub struct{}

func (appRepoStub) Search(context.Context, int) ([]domain.Application, error) {
	return []domain.Application{{ID: intPtr(1)}}, nil
}
func (appRepoStub) SearchByValue(context.Context, int, string) ([]domain.Application, error) {
	return nil, nil
}
func (appRepoStub) GetByID(context.Context, int, int) (*domain.Application, error) { return nil, nil }
func (appRepoStub) ListVersions(context.Context, int, int) ([]domain.ApplicationVersion, error) {
	return nil, nil
}
func (appRepoStub) SaveAndroid(context.Context, int, domain.Application) (*domain.Application, error) {
	return nil, nil
}
func (appRepoStub) SaveWeb(context.Context, int, domain.Application) (*domain.Application, error) {
	return nil, nil
}
func (appRepoStub) SaveVersion(context.Context, int, domain.ApplicationVersion) (*domain.ApplicationVersion, error) {
	return nil, nil
}
func (appRepoStub) DeleteApp(context.Context, int, int) error { return nil }
func (appRepoStub) DeleteVersion(context.Context, int, int) error { return nil }
func (appRepoStub) ValidatePkg(context.Context, int, domain.ValidatePkgRequest) ([]domain.Application, error) {
	return nil, nil
}
func (appRepoStub) GetAppConfigurations(context.Context, int, int) ([]domain.ApplicationConfigurationLink, error) {
	return nil, nil
}
func (appRepoStub) UpdateAppConfigurations(context.Context, int, domain.LinkConfigurationsToAppRequest) error {
	return nil
}
func (appRepoStub) GetVersionConfigurations(context.Context, int, int) ([]domain.ApplicationVersionConfigurationLink, error) {
	return nil, nil
}
func (appRepoStub) UpdateVersionConfigurations(context.Context, int, domain.LinkConfigurationsToAppVersionRequest) error {
	return nil
}
func (appRepoStub) AdminSearch(context.Context, string) ([]domain.Application, error) { return nil, nil }
func (appRepoStub) TurnIntoCommon(context.Context, int) error { return nil }
func (appRepoStub) CustomerFilesDir(context.Context, int) (string, error) {
	return "customer-1", nil
}

func intPtr(n int) *int { return &n }

func TestSearch_permissionDenied(t *testing.T) {
	svc := NewService(appRepoStub{}, nil, "")
	_, err := svc.Search(context.Background(), &platformauth.Principal{AuthLoaded: true})
	if err != ErrPermissionDenied {
		t.Fatalf("got %v", err)
	}
}

func TestSearch_ok(t *testing.T) {
	svc := NewService(appRepoStub{}, nil, "")
	p := &platformauth.Principal{AuthLoaded: true, Permissions: []string{platformauth.PermApplications}}
	list, err := svc.Search(context.Background(), p)
	if err != nil || len(list) != 1 {
		t.Fatalf("err=%v len=%d", err, len(list))
	}
}
