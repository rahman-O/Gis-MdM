package application

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type cfgRepoStub struct{}

func (cfgRepoStub) ListByCustomer(context.Context, int) ([]domain.LookupItem, error) { return nil, nil }
func (cfgRepoStub) Search(context.Context, int) ([]domain.Configuration, error) {
	return []domain.Configuration{{ID: intPtr(1)}}, nil
}
func (cfgRepoStub) SearchByValue(context.Context, int, string) ([]domain.Configuration, error) {
	return nil, nil
}
func (cfgRepoStub) GetByID(context.Context, int, int) (*domain.Configuration, error) { return nil, nil }
func (cfgRepoStub) GetByName(context.Context, int, string) (*domain.Configuration, error) {
	return nil, nil
}
func (cfgRepoStub) CountDevicesUsing(context.Context, int) (int64, error) { return 0, nil }
func (cfgRepoStub) Insert(context.Context, int, domain.Configuration) (int, error) { return 1, nil }
func (cfgRepoStub) Update(context.Context, int, domain.Configuration) error { return nil }
func (cfgRepoStub) Delete(context.Context, int, int) error                     { return nil }
func (cfgRepoStub) Copy(context.Context, int, domain.CopyRequest) (int, error) { return 2, nil }
func (cfgRepoStub) ListAllApplicationsForPicker(context.Context, int) ([]domain.ConfigurationApplication, error) {
	return nil, nil
}
func (cfgRepoStub) ListConfigurationApplications(context.Context, int, int) ([]domain.ConfigurationApplication, error) {
	return nil, nil
}
func (cfgRepoStub) UpgradeApplication(context.Context, int, domain.UpgradeApplicationRequest) error {
	return nil
}

func intPtr(n int) *int { return &n }

func TestSearch_permissionDenied(t *testing.T) {
	svc := NewService(cfgRepoStub{}, nil)
	_, err := svc.Search(context.Background(), &platformauth.Principal{AuthLoaded: true})
	if err != ErrPermissionDenied {
		t.Fatalf("got %v", err)
	}
}

func TestSearch_ok(t *testing.T) {
	svc := NewService(cfgRepoStub{}, nil)
	p := &platformauth.Principal{AuthLoaded: true, Permissions: []string{platformauth.PermConfigurations}}
	list, err := svc.Search(context.Background(), p)
	if err != nil || len(list) != 1 {
		t.Fatalf("err=%v len=%d", err, len(list))
	}
}
