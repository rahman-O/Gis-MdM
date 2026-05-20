package application

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/groups/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type stubGroupRepo struct {
	groups      []domain.Group
	deviceCount int64
}

func (s *stubGroupRepo) ListByCustomer(_ context.Context, customerID int) ([]domain.Group, error) {
	if customerID != 1 {
		return []domain.Group{}, nil
	}
	return s.groups, nil
}
func (s *stubGroupRepo) ListByValue(_ context.Context, customerID int, _ string) ([]domain.Group, error) {
	return s.ListByCustomer(context.Background(), customerID)
}
func (s *stubGroupRepo) GetByName(_ context.Context, _ int, name string) (*domain.Group, error) {
	for _, g := range s.groups {
		if g.Name == name {
			return &g, nil
		}
	}
	return nil, nil
}
func (s *stubGroupRepo) CountDevicesInGroup(_ context.Context, _ int) (int64, error) {
	return s.deviceCount, nil
}
func (s *stubGroupRepo) Insert(_ context.Context, _ int, name string) (int, error) {
	s.groups = append(s.groups, domain.Group{ID: len(s.groups) + 1, Name: name})
	return len(s.groups), nil
}
func (s *stubGroupRepo) Update(context.Context, int, domain.Group) error { return nil }
func (s *stubGroupRepo) Delete(context.Context, int, int) error         { return nil }
func (s *stubGroupRepo) GrantCreatorAccess(context.Context, int64, int) error {
	return nil
}
func (s *stubGroupRepo) UserHasAllDevices(context.Context, int64) (bool, error) { return true, nil }

func testPrincipal() *platformauth.Principal {
	return &platformauth.Principal{ID: 1, CustomerID: 1, AuthLoaded: true, Permissions: []string{"settings"}}
}

func TestList_byCustomer(t *testing.T) {
	svc := NewService(&stubGroupRepo{groups: []domain.Group{{ID: 1, Name: "General"}}})
	out, err := svc.List(context.Background(), testPrincipal())
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("len %d", len(out))
	}
}

func TestDelete_notEmpty(t *testing.T) {
	svc := NewService(&stubGroupRepo{deviceCount: 2})
	err := svc.Delete(context.Background(), testPrincipal(), 1)
	if err != ErrNotEmptyGroup {
		t.Fatalf("got %v", err)
	}
}

func TestSave_duplicateName(t *testing.T) {
	svc := NewService(&stubGroupRepo{groups: []domain.Group{{ID: 1, Name: "General"}}})
	err := svc.Save(context.Background(), testPrincipal(), domain.Group{Name: "General"})
	if err != ErrDuplicateGroup {
		t.Fatalf("got %v", err)
	}
}
