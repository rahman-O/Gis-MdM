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
}
