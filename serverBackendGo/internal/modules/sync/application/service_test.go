package application

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/sync/domain"
)

type stubRepo struct{}

func (stubRepo) FindByNumber(ctx context.Context, number string) (*domain.DeviceRecord, error) {
	return &domain.DeviceRecord{ID: 1, Number: number, ConfigurationID: 1, CustomerID: 1}, nil
}
func (stubRepo) FindByOldNumber(context.Context, string) (*domain.DeviceRecord, error) {
	return nil, nil
}
func (stubRepo) CreateOnDemand(context.Context, string, domain.DeviceCreateOptions, int64) (*domain.DeviceRecord, error) {
	return nil, nil
}
func (stubRepo) CompleteMigration(context.Context, int64) error  { return nil }
func (stubRepo) TouchLastUpdate(context.Context, int64) error    { return nil }
func (stubRepo) UpdateInfo(context.Context, int64, string, string) error {
	return nil
}
func (stubRepo) UpdateCustomProps(context.Context, int64, *string, *string, *string) error {
	return nil
}
func (stubRepo) SaveApplicationSettings(context.Context, int64, []domain.SyncApplicationSetting) error {
	return nil
}
func (stubRepo) BuildSyncResponse(context.Context, domain.DeviceRecord, string, string, string, string, string) (*domain.SyncResponse, error) {
	return &domain.SyncResponse{DeviceID: "hmdm-001"}, nil
}
func (stubRepo) CountCustomers(context.Context) (int, error) { return 1, nil }

func TestGetConfiguration_secureEnrollment(t *testing.T) {
	svc := NewService(stubRepo{}, Config{SecureEnrollment: true, HashSecret: "s"})
	_, err := svc.GetConfiguration(context.Background(), "d1", "bad", "")
	if err != ErrPermissionDenied {
		t.Fatalf("want permission denied, got %v", err)
	}
}
