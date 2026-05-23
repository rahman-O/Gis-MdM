package application

import (
	"context"
	"encoding/json"

	"github.com/gis-mdm/server-backend-go/internal/modules/settings/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/settings/port"
)

// Service implements settings use cases.
type Service struct {
	repo port.SettingsRepository
}

func NewService(repo port.SettingsRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, customerID int) (*domain.Settings, error) {
	return s.repo.GetByCustomerID(ctx, customerID)
}

func (s *Service) MergeAndSaveMisc(ctx context.Context, customerID int, body map[string]interface{}) error {
	cur, err := s.repo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return err
	}
	applyMap(cur, body)
	cur.CustomerID = customerID
	if single, _ := s.repo.IsSingleCustomer(ctx); !single {
		cur.CreateNewDevices = false
		cur.NewDeviceConfigurationID = nil
	}
	return s.repo.SaveMisc(ctx, cur)
}

func (s *Service) MergeAndSaveLang(ctx context.Context, customerID int, body map[string]interface{}) error {
	cur, err := s.repo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return err
	}
	applyMap(cur, body)
	cur.CustomerID = customerID
	return s.repo.SaveLanguage(ctx, cur)
}

func (s *Service) GetUserRoleSettings(ctx context.Context, customerID, roleID int) (*domain.UserRoleSettings, error) {
	return s.repo.GetUserRoleSettings(ctx, customerID, roleID)
}

func (s *Service) SaveUserRoleSettings(ctx context.Context, customerID int, body map[string]interface{}) error {
	b, _ := json.Marshal(body)
	var settings domain.UserRoleSettings
	if err := json.Unmarshal(b, &settings); err != nil {
		return err
	}
	if rid, ok := body["roleId"].(float64); ok {
		settings.RoleID = int(rid)
	}
	settings.CustomerID = customerID
	return s.repo.SaveUserRoleSettings(ctx, customerID, settings)
}

func (s *Service) MergeAndSaveDesign(ctx context.Context, customerID int, body map[string]interface{}) error {
	cur, err := s.repo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return err
	}
	applyMap(cur, body)
	cur.CustomerID = customerID
	return s.repo.SaveDesign(ctx, cur)
}

func applyMap(s *domain.Settings, body map[string]interface{}) {
	b, _ := json.Marshal(body)
	var patch domain.Settings
	_ = json.Unmarshal(b, &patch)
	if v, ok := body["language"].(string); ok {
		s.Language = v
	}
	if v, ok := body["useDefaultLanguage"].(bool); ok {
		s.UseDefaultLanguage = v
	}
	if v, ok := body["createNewDevices"].(bool); ok {
		s.CreateNewDevices = v
	}
	if v, ok := body["passwordLength"].(float64); ok {
		s.PasswordLength = int(v)
	}
	if v, ok := body["passwordStrength"].(float64); ok {
		s.PasswordStrength = int(v)
	}
	if v, ok := body["twoFactor"].(bool); ok {
		s.TwoFactor = v
	}
	if v, ok := body["idleLogout"].(float64); ok {
		n := int(v)
		s.IdleLogout = &n
	} else if body["idleLogout"] == nil {
		s.IdleLogout = nil
	}
	if patch.BackgroundColor != "" {
		s.BackgroundColor = patch.BackgroundColor
	}
	if patch.TextColor != "" {
		s.TextColor = patch.TextColor
	}
	if patch.BackgroundImageURL != "" {
		s.BackgroundImageURL = patch.BackgroundImageURL
	}
	if patch.IconSize != "" {
		s.IconSize = patch.IconSize
	}
	if patch.DesktopHeader != "" {
		s.DesktopHeader = patch.DesktopHeader
	}
	if v, ok := body["newDeviceGroupId"].(float64); ok {
		n := int(v)
		s.NewDeviceGroupID = &n
	} else if body["newDeviceGroupId"] == nil {
		s.NewDeviceGroupID = nil
	}
	if v, ok := body["phoneNumberFormat"].(string); ok {
		s.PhoneNumberFormat = v
	}
	if v, ok := body["customPropertyName1"].(string); ok {
		s.CustomPropertyName1 = v
	}
	if v, ok := body["customPropertyName2"].(string); ok {
		s.CustomPropertyName2 = v
	}
	if v, ok := body["customPropertyName3"].(string); ok {
		s.CustomPropertyName3 = v
	}
	if v, ok := body["customMultiline1"].(bool); ok {
		s.CustomMultiline1 = v
	}
	if v, ok := body["customMultiline2"].(bool); ok {
		s.CustomMultiline2 = v
	}
	if v, ok := body["customMultiline3"].(bool); ok {
		s.CustomMultiline3 = v
	}
	if v, ok := body["customSend1"].(bool); ok {
		s.CustomSend1 = v
	}
	if v, ok := body["customSend2"].(bool); ok {
		s.CustomSend2 = v
	}
	if v, ok := body["customSend3"].(bool); ok {
		s.CustomSend3 = v
	}
	if v, ok := body["desktopHeaderTemplate"].(string); ok {
		s.DesktopHeaderTemplate = v
	}
	if v, ok := body["sendDescription"].(bool); ok {
		s.SendDescription = v
	}
}
