package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	authdomain "github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/customers/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/customers/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
	"github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

// Service implements customer use cases.
type Service struct {
	customers port.CustomerRepository
	users     port.UserLookup
}

func NewService(customers port.CustomerRepository, users port.UserLookup) *Service {
	return &Service{customers: customers, users: users}
}

var (
	ErrPermissionDenied      = errors.New("error.permission.denied")
	ErrDuplicateCustomerName = errors.New("error.duplicate.customer.name")
	ErrDuplicateEmail        = errors.New("error.duplicate.email")
	ErrOrgAdminNotFound      = errors.New("error.notfound.customer.admin")
	ErrImpersonationBlocked  = errors.New("impersonation.blocked")
)

func (s *Service) requireSuperAdmin(p *platformauth.Principal) error {
	if p == nil || !p.SuperAdmin {
		return ErrPermissionDenied
	}
	return nil
}

func (s *Service) Search(ctx context.Context, p *platformauth.Principal, req domain.SearchRequest) (*domain.Paginated, error) {
	if err := s.requireSuperAdmin(p); err != nil {
		return nil, err
	}
	req.NormalizeSearch()
	items, err := s.customers.Search(ctx, req)
	if err != nil {
		return nil, err
	}
	total, err := s.customers.Count(ctx, req)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.Customer{}
	}
	return &domain.Paginated{Items: items, TotalItemsCount: total}, nil
}

func (s *Service) GetForEdit(ctx context.Context, p *platformauth.Principal, id int) (*domain.Customer, error) {
	if err := s.requireSuperAdmin(p); err != nil {
		return nil, err
	}
	return s.customers.GetByID(ctx, id)
}

func (s *Service) PrefixUsed(ctx context.Context, p *platformauth.Principal, prefix string) (bool, error) {
	if err := s.requireSuperAdmin(p); err != nil {
		return false, err
	}
	return s.customers.PrefixUsed(ctx, prefix)
}

func (s *Service) Delete(ctx context.Context, p *platformauth.Principal, id int) error {
	if err := s.requireSuperAdmin(p); err != nil {
		return err
	}
	return s.customers.Delete(ctx, id)
}

// Save creates or updates a customer. Default devices/config copy deferred until Phase 4/5.
func (s *Service) Save(ctx context.Context, p *platformauth.Principal, c domain.Customer) (map[string]string, error) {
	if err := s.requireSuperAdmin(p); err != nil {
		return nil, err
	}
	email := strings.TrimSpace(c.Email)
	if c.ID == nil {
		return s.createCustomer(ctx, c, email)
	}
	if err := s.validateDuplicates(ctx, c, email); err != nil {
		return nil, err
	}
	if err := s.customers.Update(ctx, &c); err != nil {
		return nil, err
	}
	if err := s.syncOrgAdmin(ctx, *c.ID, c); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Service) createCustomer(ctx context.Context, c domain.Customer, email string) (map[string]string, error) {
	if err := s.validateDuplicates(ctx, c, email); err != nil {
		return nil, err
	}
	c.FilesDir = fmt.Sprintf("%d", time.Now().UnixNano())
	id, err := s.customers.Insert(ctx, &c)
	if err != nil {
		return nil, err
	}
	plain := crypto.GeneratePassword(8)
	hash := crypto.HashFromMd5(crypto.MD5UpperHex(plain))
	login := domain.Transliterate(c.Name)
	name := c.MainUserName
	if name == "" {
		name = c.Name
	}
	token := crypto.GenerateAuthToken()
	if err := s.users.InsertOrgAdmin(ctx, id, login, name, email, hash, token, false); err != nil {
		return nil, err
	}
	return map[string]string{"adminCredentials": login + "/" + plain}, nil
}

func (s *Service) syncOrgAdmin(ctx context.Context, customerID int, c domain.Customer) error {
	admin, err := s.users.FindOrgAdmin(ctx, customerID)
	if err != nil || admin == nil {
		return nil
	}
	login := domain.Transliterate(c.Name)
	name := c.MainUserName
	if name == "" {
		name = c.Name
	}
	return s.users.UpdateOrgAdminMainDetails(ctx, admin.ID, login, name, strings.TrimSpace(c.Email))
}

func (s *Service) validateDuplicates(ctx context.Context, c domain.Customer, email string) error {
	if byName, _ := s.customers.GetByName(ctx, c.Name); byName != nil {
		if c.ID == nil || byName.ID == nil || *byName.ID != *c.ID {
			return ErrDuplicateCustomerName
		}
	}
	if email != "" {
		if byEmail, _ := s.customers.GetByEmail(ctx, email); byEmail != nil {
			if c.ID == nil || byEmail.ID == nil || *byEmail.ID != *c.ID {
				return ErrDuplicateEmail
			}
		}
		if u, _ := s.users.FindByEmail(ctx, email); u != nil {
			if c.ID == nil || u.CustomerID != *c.ID {
				return ErrDuplicateEmail
			}
		}
	}
	if c.ID == nil {
		if u, _ := s.users.FindByLogin(ctx, domain.Transliterate(c.Name)); u != nil {
			return ErrDuplicateCustomerName
		}
	}
	return nil
}

func (s *Service) Impersonate(ctx context.Context, p *platformauth.Principal, customerID int) (*authdomain.UserView, error) {
	if err := s.requireSuperAdmin(p); err != nil {
		return nil, err
	}
	admin, err := s.users.FindOrgAdmin(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrOrgAdminNotFound
	}
	token, err := s.users.EnsureAuthToken(ctx, admin.ID)
	if err != nil {
		if errors.Is(err, port.ErrImpersonationBlocked) {
			return nil, ErrImpersonationBlocked
		}
		return nil, err
	}
	admin.AuthToken = token
	admin.Password = ""
	if admin.UserRole != nil {
		admin.UserRole.SuperAdmin = false
	}
	v := authdomain.NewUserView(admin)
	return v, nil
}
