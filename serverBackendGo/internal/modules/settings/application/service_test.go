package application

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/settings/domain"
)

type stubSettingsRepo struct {
	userRole *domain.UserRoleSettings
}

func (s *stubSettingsRepo) GetByCustomerID(context.Context, int) (*domain.Settings, error) {
	return &domain.Settings{}, nil
}
func (s *stubSettingsRepo) IsSingleCustomer(context.Context) (bool, error) { return true, nil }
func (s *stubSettingsRepo) SaveMisc(context.Context, *domain.Settings) error   { return nil }
func (s *stubSettingsRepo) SaveLanguage(context.Context, *domain.Settings) error { return nil }
func (s *stubSettingsRepo) SaveDesign(context.Context, *domain.Settings) error { return nil }
func (s *stubSettingsRepo) GetUserRoleSettings(context.Context, int, int) (*domain.UserRoleSettings, error) {
	if s.userRole != nil {
		return s.userRole, nil
	}
	def := domain.DefaultUserRoleSettings(2, 1)
	return &def, nil
}
func (s *stubSettingsRepo) SaveUserRoleSettings(context.Context, int, domain.UserRoleSettings) error {
	return nil
}

func TestGetUserRoleSettings_defaultsWhenMissing(t *testing.T) {
	svc := NewService(&stubSettingsRepo{})
	got, err := svc.GetUserRoleSettings(context.Background(), 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if got.RoleID != 2 || !got.ColumnDisplayedDeviceNumber {
		t.Fatalf("unexpected %+v", got)
	}
}
