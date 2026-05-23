package application

import (
	"context"
	"testing"

	authdomain "github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/customers/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/customers/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type stubRepo struct {
	items []domain.Customer
	total int64
}

func (s stubRepo) Search(context.Context, domain.SearchRequest) ([]domain.Customer, error) {
	return s.items, nil
}
func (s stubRepo) Count(context.Context, domain.SearchRequest) (int64, error) {
	return s.total, nil
}
func (s stubRepo) GetByID(context.Context, int) (*domain.Customer, error) { return nil, nil }
func (s stubRepo) GetByName(context.Context, string) (*domain.Customer, error) {
	return nil, nil
}
func (s stubRepo) GetByEmail(context.Context, string) (*domain.Customer, error) {
	return nil, nil
}
func (s stubRepo) Insert(context.Context, *domain.Customer) (int, error) { return 1, nil }
func (s stubRepo) Update(context.Context, *domain.Customer) error          { return nil }
func (s stubRepo) Delete(context.Context, int) error                     { return nil }
func (s stubRepo) PrefixUsed(context.Context, string) (bool, error)      { return false, nil }

type stubUsers struct {
	admin *authdomain.User
}

func (s stubUsers) FindOrgAdmin(context.Context, int) (*authdomain.User, error) {
	return s.admin, nil
}
func (s stubUsers) FindByLogin(context.Context, string) (*authdomain.User, error) {
	return nil, nil
}
func (s stubUsers) FindByEmail(context.Context, string) (*authdomain.User, error) {
	return nil, nil
}
func (s stubUsers) EnsureAuthToken(context.Context, int64) (string, error) {
	return "tok", nil
}
func (s stubUsers) UpdateOrgAdminMainDetails(context.Context, int64, string, string, string) error {
	return nil
}
func (s stubUsers) InsertOrgAdmin(context.Context, int, string, string, string, string, string, bool) error {
	return nil
}

func TestSearch_requiresSuperAdmin(t *testing.T) {
	svc := NewService(stubRepo{}, stubUsers{})
	_, err := svc.Search(context.Background(), &platformauth.Principal{SuperAdmin: false}, domain.SearchRequest{})
	if err != ErrPermissionDenied {
		t.Fatalf("want permission denied, got %v", err)
	}
}

func TestSearch_ok(t *testing.T) {
	svc := NewService(stubRepo{items: []domain.Customer{{Name: "A"}}, total: 1}, stubUsers{})
	p := &platformauth.Principal{SuperAdmin: true}
	out, err := svc.Search(context.Background(), p, domain.SearchRequest{CurrentPage: 1, PageSize: 10})
	if err != nil || out.TotalItemsCount != 1 || len(out.Items) != 1 {
		t.Fatalf("search: %v %+v", err, out)
	}
}

func TestImpersonate_noAdmin(t *testing.T) {
	svc := NewService(stubRepo{}, stubUsers{admin: nil})
	p := &platformauth.Principal{SuperAdmin: true}
	_, err := svc.Impersonate(context.Background(), p, 2)
	if err != ErrOrgAdminNotFound {
		t.Fatalf("want not found admin, got %v", err)
	}
}

func TestSave_duplicateName(t *testing.T) {
	id := 2
	svc := NewService(&dupRepo{id: id}, stubUsers{})
	p := &platformauth.Principal{SuperAdmin: true}
	_, err := svc.Save(context.Background(), p, domain.Customer{Name: "Taken"})
	if err != ErrDuplicateCustomerName {
		t.Fatalf("want duplicate name, got %v", err)
	}
}

type dupRepo struct {
	id int
	stubRepo
}

func (d *dupRepo) GetByName(context.Context, string) (*domain.Customer, error) {
	return &domain.Customer{ID: &d.id, Name: "Taken"}, nil
}

func TestImpersonate_blockedResetToken(t *testing.T) {
	admin := &authdomain.User{ID: 10, Login: "org", UserRole: &authdomain.UserRole{ID: 2}}
	svc := NewService(stubRepo{}, blockedUsers{admin: admin})
	p := &platformauth.Principal{SuperAdmin: true}
	_, err := svc.Impersonate(context.Background(), p, 1)
	if err != ErrImpersonationBlocked {
		t.Fatalf("want blocked, got %v", err)
	}
}

type blockedUsers struct {
	stubUsers
	admin *authdomain.User
}

func (b blockedUsers) FindOrgAdmin(context.Context, int) (*authdomain.User, error) {
	return b.admin, nil
}
func (blockedUsers) EnsureAuthToken(context.Context, int64) (string, error) {
	return "", port.ErrImpersonationBlocked
}
