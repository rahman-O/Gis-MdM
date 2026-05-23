package application

import (
	"context"
	"testing"

	authdomain "github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/users/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type stubRepo struct {
	users map[int64]*authdomain.User
	list  []*authdomain.User
}

func (s *stubRepo) FindByID(_ context.Context, id int64) (*authdomain.User, error) {
	return s.users[id], nil
}
func (s *stubRepo) FindByLogin(_ context.Context, login string) (*authdomain.User, error) {
	for _, u := range s.users {
		if u.Login == login {
			return u, nil
		}
	}
	return nil, nil
}
func (s *stubRepo) FindByEmail(_ context.Context, email string) (*authdomain.User, error) {
	for _, u := range s.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}
func (s *stubRepo) ListByCustomer(context.Context, int, string) ([]*authdomain.User, error) {
	if s.list != nil {
		return s.list, nil
	}
	return nil, nil
}
func (s *stubRepo) UpdateMainDetails(context.Context, *authdomain.User) error { return nil }
func (s *stubRepo) UpdatePassword(context.Context, int64, string, string, bool, bool, *string) error {
	return nil
}
func (s *stubRepo) Insert(context.Context, *authdomain.User) error { return nil }
func (s *stubRepo) Delete(context.Context, int64) error              { return nil }
func (s *stubRepo) IsSingleCustomer(context.Context) (bool, error)   { return true, nil }
func (s *stubRepo) GetCustomerSettings(context.Context, int) (*authdomain.CustomerSettings, error) {
	return &authdomain.CustomerSettings{}, nil
}
func (s *stubRepo) PasswordResetEnabled(context.Context, int) (bool, error) { return false, nil }

func TestUpdateProfile_duplicateEmail(t *testing.T) {
	repo := &stubRepo{users: map[int64]*authdomain.User{
		1: {ID: 1, Email: "a@test.local", Name: "A"},
		2: {ID: 2, Email: "b@test.local", Name: "B"},
	}}
	svc := NewService(repo)
	p := &platformauth.Principal{ID: 1}
	_, err := svc.UpdateProfile(context.Background(), p, domain.ProfilePayload{
		ID: 1, Name: "A", Email: "b@test.local",
	})
	if err != ErrDuplicateEmail {
		t.Fatalf("want duplicate email, got %v", err)
	}
}

func TestChangePassword_wrongOld(t *testing.T) {
	repo := &stubRepo{users: map[int64]*authdomain.User{
		1: {ID: 1, Login: "u", Password: "HASH"},
	}}
	svc := NewService(repo)
	id := int64(1)
	p := &platformauth.Principal{ID: 1}
	err := svc.ChangePassword(context.Background(), p, domain.UserPayload{
		ID: &id, Login: "u", OldPassword: "BAD", NewPassword: "NEWMD5",
	})
	if err != ErrWrongPassword {
		t.Fatalf("want wrong password, got %v", err)
	}
}
