package application

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/devices/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/devices/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type stubDeviceRepo struct {
	total int64
}

func (s *stubDeviceRepo) LoadUserScope(_ context.Context, _ int64) (*port.UserScope, error) {
	return &port.UserScope{UserID: 1, CustomerID: 1, AllDevicesAvailable: true}, nil
}
func (s *stubDeviceRepo) Search(_ context.Context, _ port.UserScope, req domain.SearchRequest) ([]domain.DeviceView, error) {
	return []domain.DeviceView{{ID: 1, Number: "hmdm-001"}}, nil
}
func (s *stubDeviceRepo) Count(_ context.Context, _ port.UserScope, _ domain.SearchRequest) (int64, error) {
	return s.total, nil
}
func (s *stubDeviceRepo) ListConfigurations(_ context.Context, _ int) (map[int]domain.ConfigurationView, error) {
	name := "Default"
	return map[int]domain.ConfigurationView{1: {ID: 1, Name: &name}}, nil
}
func (s *stubDeviceRepo) GetByNumber(context.Context, port.UserScope, string) (*domain.DeviceView, error) {
	return nil, nil
}
func (s *stubDeviceRepo) GetByID(context.Context, int, int) (*domain.DeviceView, error) { return nil, nil }
func (s *stubDeviceRepo) ExistsNumber(context.Context, int, string, int) (bool, error)     { return false, nil }
func (s *stubDeviceRepo) CountDevices(context.Context, int) (int64, error)                { return 0, nil }
func (s *stubDeviceRepo) DeviceLimit(context.Context, int) (int, error)                   { return 0, nil }
func (s *stubDeviceRepo) Insert(context.Context, int, domain.SaveDevice) (int, error)   { return 1, nil }
func (s *stubDeviceRepo) Update(context.Context, int, domain.SaveDevice) error          { return nil }
func (s *stubDeviceRepo) UpdateConfigurationBulk(context.Context, int, []int, int) error {
	return nil
}
func (s *stubDeviceRepo) Delete(context.Context, int, int) error { return nil }
func (s *stubDeviceRepo) DeleteBulk(context.Context, int, []int) error {
	return nil
}
func (s *stubDeviceRepo) UpdateGroupBulk(context.Context, int, domain.GroupBulkRequest) error {
	return nil
}
func (s *stubDeviceRepo) Autocomplete(context.Context, port.UserScope, string, int) ([]domain.LookupItem, error) {
	return nil, nil
}
func (s *stubDeviceRepo) UpdateDescription(context.Context, int, int, string) error { return nil }
func (s *stubDeviceRepo) ListAppSettings(context.Context, int) ([]domain.AppSetting, error) {
	return nil, nil
}
func (s *stubDeviceRepo) SaveAppSettings(context.Context, int, []domain.AppSetting) error { return nil }

func devicePrincipal(perms ...string) *platformauth.Principal {
	return &platformauth.Principal{ID: 1, CustomerID: 1, AuthLoaded: true, Permissions: perms}
}

func TestSearch_pagination(t *testing.T) {
	svc := NewService(&stubDeviceRepo{total: 5}, nil)
	req := domain.SearchRequest{PageNum: 1, PageSize: 10}
	out, err := svc.Search(context.Background(), devicePrincipal(), req)
	if err != nil {
		t.Fatal(err)
	}
	if out.Devices.TotalItemsCount != 5 {
		t.Fatalf("total %d", out.Devices.TotalItemsCount)
	}
	if len(out.Devices.Items) != 1 {
		t.Fatalf("items %d", len(out.Devices.Items))
	}
}

func TestSearchRequest_prepareWrapsValue(t *testing.T) {
	v := "abc"
	req := domain.SearchRequest{PageNum: 1, PageSize: 10, Value: &v}
	req.Prepare()
	if req.Value == nil || *req.Value != "%abc%" {
		t.Fatalf("value %v", req.Value)
	}
}

func TestDelete_permissionDenied(t *testing.T) {
	svc := NewService(&stubDeviceRepo{}, nil)
	err := svc.Delete(context.Background(), devicePrincipal(), 1)
	if err != ErrPermissionDenied {
		t.Fatalf("got %v", err)
	}
}
