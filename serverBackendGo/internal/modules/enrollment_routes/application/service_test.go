package application_test

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type fakeRepo struct {
	treeOK    bool
	createdID int
	kind      string
}

func (f *fakeRepo) ListViews(context.Context, int) ([]domain.EnrollmentRouteView, error) {
	return nil, nil
}

func (f *fakeRepo) GetViewByID(context.Context, int, int) (*domain.EnrollmentRouteView, error) {
	return &domain.EnrollmentRouteView{
		ID:                     1,
		QRCodeKey:              "k",
		TargetNodeID:           2,
		DeviceIdentityMode:     "imei",
		BootstrapIntent:        domain.BootstrapIntentStable,
		BootstrapApplicationID: 1,
		Status:                 "active",
	}, nil
}

func (f *fakeRepo) Create(context.Context, int, domain.CreateRequest, string, domain.ResolvedBootstrap, bool) (int, error) {
	return f.createdID, nil
}

func (f *fakeRepo) Update(context.Context, int, int, domain.UpdateRequest, *domain.ResolvedBootstrap, *bool) error {
	return nil
}

func (f *fakeRepo) Delete(context.Context, int, int) error { return nil }

func (f *fakeRepo) DeleteImpact(context.Context, int, int) (*domain.EnrollmentDeleteImpact, error) {
	return &domain.EnrollmentDeleteImpact{}, nil
}

func (f *fakeRepo) TreeNodeBelongsToCustomer(context.Context, int, int) (bool, error) {
	return f.treeOK, nil
}

func (f *fakeRepo) NodePlacementKind(context.Context, int, int) (string, error) {
	if f.kind != "" {
		return f.kind, nil
	}
	return "locked", nil
}

func (f *fakeRepo) ListTreeNodeOptions(context.Context, int, int) ([]domain.TreeNodeOption, error) {
	return nil, nil
}

func (f *fakeRepo) ListBootstrapApps(context.Context, int) ([]domain.BootstrapAppOption, error) {
	return nil, nil
}

func (f *fakeRepo) RecordQRViewed(context.Context, int) error { return nil }

func principal() *platformauth.Principal {
	return &platformauth.Principal{CustomerID: 1, Permissions: []string{"configurations"}}
}

func TestCreate_RequiresTreeNode(t *testing.T) {
	svc := application.NewService(&fakeRepo{treeOK: false}, nil, 500)
	_, err := svc.Create(context.Background(), principal(), domain.CreateRequest{
		Name:                   "R1",
		TargetNodeID:           2,
		BootstrapApplicationID: 1,
		BootstrapIntent:        domain.BootstrapIntentStable,
	})
	if err != application.ErrTreeNodeRequired {
		t.Fatalf("expected tree node error, got %v", err)
	}
}

func TestCreate_RequiresContainerAck(t *testing.T) {
	svc := application.NewService(&fakeRepo{treeOK: true, kind: "inheritable"}, nil, 500)
	_, err := svc.Create(context.Background(), principal(), domain.CreateRequest{
		Name:                   "R1",
		TargetNodeID:           2,
		BootstrapApplicationID: 1,
		BootstrapIntent:        domain.BootstrapIntentStable,
	})
	if err != application.ErrContainerAckRequired {
		t.Fatalf("expected container ack error, got %v", err)
	}
}
