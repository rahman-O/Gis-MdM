package application_test

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type fakeRepo struct {
	publishedOK bool
	treeOK      bool
	createdID   int
}

func (f *fakeRepo) List(context.Context, int) ([]domain.Route, error) { return nil, nil }

func (f *fakeRepo) GetByID(context.Context, int, int) (*domain.RouteDetail, error) {
	return &domain.RouteDetail{Route: domain.Route{ID: 1, QRCodeKey: "k"}}, nil
}

func (f *fakeRepo) Create(context.Context, int, domain.CreateRequest, string) (int, error) {
	return f.createdID, nil
}

func (f *fakeRepo) Update(context.Context, int, int, domain.UpdateRequest, *string) error {
	return nil
}

func (f *fakeRepo) IsPublishedProfileVersion(context.Context, int, int) (bool, error) {
	return f.publishedOK, nil
}

func (f *fakeRepo) ListPublishedProfileVersions(context.Context, int) ([]domain.PublishedProfileVersion, error) {
	return nil, nil
}

func (f *fakeRepo) TreeNodeBelongsToCustomer(context.Context, int, int) (bool, error) {
	return f.treeOK, nil
}

func principal() *platformauth.Principal {
	return &platformauth.Principal{CustomerID: 1, Permissions: []string{"configurations"}}
}

func intPtr(n int) *int { return &n }

func TestCreate_RequiresPublishedProfileVersion(t *testing.T) {
	svc := application.NewService(&fakeRepo{publishedOK: false, treeOK: true})
	main := 1
	_, err := svc.Create(context.Background(), principal(), domain.CreateRequest{
		Name:              "R1",
		ProfileVersionID:  intPtr(10),
		DefaultTreeNodeID: 2,
		MainAppID:         &main,
	})
	if err != application.ErrPublishedVersionRequired {
		t.Fatalf("expected published version error, got %v", err)
	}
}

func TestCreate_RequiresTreeNode(t *testing.T) {
	svc := application.NewService(&fakeRepo{publishedOK: true, treeOK: false})
	main := 1
	_, err := svc.Create(context.Background(), principal(), domain.CreateRequest{
		Name:              "R1",
		ProfileVersionID:  intPtr(10),
		DefaultTreeNodeID: 2,
		MainAppID:         &main,
	})
	if err != application.ErrTreeNodeRequired {
		t.Fatalf("expected tree node error, got %v", err)
	}
}

func TestCreate_SucceedsWhenValid(t *testing.T) {
	svc := application.NewService(&fakeRepo{publishedOK: true, treeOK: true, createdID: 5})
	main := 1
	detail, err := svc.Create(context.Background(), principal(), domain.CreateRequest{
		Name:              "R1",
		ProfileVersionID:  intPtr(10),
		DefaultTreeNodeID: 2,
		MainAppID:         &main,
	})
	if err != nil {
		t.Fatal(err)
	}
	if detail.ID != 1 {
		t.Fatalf("expected loaded detail id 1, got %d", detail.ID)
	}
}

func TestCreate_WithoutProfileVersion(t *testing.T) {
	svc := application.NewService(&fakeRepo{treeOK: true, createdID: 7})
	main := 1
	detail, err := svc.Create(context.Background(), principal(), domain.CreateRequest{
		Name:              "R2",
		DefaultTreeNodeID: 2,
		MainAppID:         &main,
	})
	if err != nil {
		t.Fatal(err)
	}
	if detail == nil {
		t.Fatal("expected detail")
	}
}
