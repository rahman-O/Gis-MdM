package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/icons/domain"
)

// IconRepository persists launcher icons.
type IconRepository interface {
	List(ctx context.Context, customerID int, filter string) ([]domain.Icon, error)
	Save(ctx context.Context, icon domain.Icon) (*domain.Icon, error)
	Delete(ctx context.Context, customerID, id int) error
}
