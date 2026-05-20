package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/settings/domain"
)

// SettingsRepository loads and updates tenant settings.
type SettingsRepository interface {
	GetByCustomerID(ctx context.Context, customerID int) (*domain.Settings, error)
	IsSingleCustomer(ctx context.Context) (bool, error)
	SaveMisc(ctx context.Context, s *domain.Settings) error
	SaveLanguage(ctx context.Context, s *domain.Settings) error
	SaveDesign(ctx context.Context, s *domain.Settings) error
}
