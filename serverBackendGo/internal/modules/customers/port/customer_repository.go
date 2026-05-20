package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/customers/domain"
)

// CustomerRepository persists tenant customers.
type CustomerRepository interface {
	Search(ctx context.Context, req domain.SearchRequest) ([]domain.Customer, error)
	Count(ctx context.Context, req domain.SearchRequest) (int64, error)
	GetByID(ctx context.Context, id int) (*domain.Customer, error)
	GetByName(ctx context.Context, name string) (*domain.Customer, error)
	GetByEmail(ctx context.Context, email string) (*domain.Customer, error)
	Insert(ctx context.Context, c *domain.Customer) (int, error)
	Update(ctx context.Context, c *domain.Customer) error
	Delete(ctx context.Context, id int) error
	PrefixUsed(ctx context.Context, prefix string) (bool, error)
}
