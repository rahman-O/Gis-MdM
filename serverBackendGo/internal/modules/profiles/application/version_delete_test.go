package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
)

type deleteRepo struct {
	fakeRepo
	eligibility [3]bool
	deleteErr   error
}

func (d *deleteRepo) VersionDeleteEligibility(context.Context, int, int, int) (bool, bool, bool, error) {
	return d.eligibility[0], d.eligibility[1], d.eligibility[2], nil
}

func (d *deleteRepo) DeleteVersion(context.Context, int, int, int) error {
	return d.deleteErr
}

func TestVersionDelete_ActivePublishedRejected(t *testing.T) {
	repo := &deleteRepo{
		fakeRepo:    fakeRepo{meta: &domain.ProfileMeta{ID: 1, Name: "P"}},
		eligibility: [3]bool{true, false, false},
	}
	svc := application.NewVersionDeleteService(repo, nil)
	_, err := svc.Delete(context.Background(), principal(), 1, 5)
	if !errors.Is(err, application.ErrVersionDeleteActivePublished) {
		t.Fatalf("expected active published error, got %v", err)
	}
}

func TestVersionDelete_AssignedRejected(t *testing.T) {
	repo := &deleteRepo{
		fakeRepo:    fakeRepo{meta: &domain.ProfileMeta{ID: 1, Name: "P"}},
		eligibility: [3]bool{false, true, false},
	}
	svc := application.NewVersionDeleteService(repo, nil)
	_, err := svc.Delete(context.Background(), principal(), 1, 5)
	if !errors.Is(err, application.ErrVersionDeleteAssigned) {
		t.Fatalf("expected assigned error, got %v", err)
	}
}

func TestVersionDelete_Success(t *testing.T) {
	repo := &deleteRepo{
		fakeRepo: fakeRepo{meta: &domain.ProfileMeta{ID: 1, Name: "P"}},
	}
	svc := application.NewVersionDeleteService(repo, nil)
	res, err := svc.Delete(context.Background(), principal(), 1, 9)
	if err != nil {
		t.Fatal(err)
	}
	if res.VersionID != 9 {
		t.Fatalf("unexpected result %+v", res)
	}
}

func TestVersionDelete_MapsPostgresNotFound(t *testing.T) {
	repo := &deleteRepo{
		fakeRepo:  fakeRepo{meta: &domain.ProfileMeta{ID: 1, Name: "P"}},
		deleteErr: postgres.ErrVersionNotFound,
	}
	svc := application.NewVersionDeleteService(repo, nil)
	_, err := svc.Delete(context.Background(), principal(), 1, 9)
	if !errors.Is(err, application.ErrVersionNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}
