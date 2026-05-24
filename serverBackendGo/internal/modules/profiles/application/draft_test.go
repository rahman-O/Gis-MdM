package application_test

import (
	"context"
	"testing"

	cfgdomain "github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type fakeRepo struct {
	meta           *domain.ProfileMeta
	ensureDraftID  int
	ensureDraftErr error
	saveErr        error
}

func (f *fakeRepo) List(context.Context, int) ([]domain.ProfileListItem, error) { return nil, nil }

func (f *fakeRepo) GetMeta(context.Context, int, int) (*domain.ProfileMeta, error) {
	return f.meta, nil
}

func (f *fakeRepo) GetVersion(context.Context, int, int, int) (*cfgdomain.Configuration, *domain.VersionMeta, error) {
	return nil, nil, nil
}

func (f *fakeRepo) Create(context.Context, int, domain.CreateRequest) (int, int, error) {
	return 1, 10, nil
}

func (f *fakeRepo) EnsureDraft(context.Context, int, int) (int, error) {
	return f.ensureDraftID, f.ensureDraftErr
}

func (f *fakeRepo) SaveDraft(context.Context, int, int, int, cfgdomain.Configuration) error {
	return f.saveErr
}

func (f *fakeRepo) ListVersionApplications(context.Context, int, int) ([]cfgdomain.ConfigurationApplication, error) {
	return nil, nil
}

func (f *fakeRepo) CountImpact(context.Context, int, int) (int, int, error) { return 0, 0, nil }

func (f *fakeRepo) PublishVersion(context.Context, int, int, int, []byte, string, int) error {
	return nil
}

func (f *fakeRepo) InsertDomainEvent(context.Context, string, string, []byte) error { return nil }

func (f *fakeRepo) ListVersions(context.Context, int, int) ([]domain.VersionListItem, error) {
	return nil, nil
}

func (f *fakeRepo) ForkDraftFromPublished(context.Context, int, int, int) (int, error) {
	return 11, nil
}

func (f *fakeRepo) DeleteVersion(context.Context, int, int, int) error { return nil }

func (f *fakeRepo) VersionDeleteEligibility(context.Context, int, int, int) (bool, bool, bool, error) {
	return false, false, false, nil
}

func principal() *platformauth.Principal {
	return &platformauth.Principal{CustomerID: 1, Permissions: []string{"configurations"}}
}

func TestGetMeta_AutoEnsuresDraftWhenMissing(t *testing.T) {
	repo := &fakeRepo{
		meta:          &domain.ProfileMeta{ID: 1, Name: "P", PublishedVersionID: intPtr(5)},
		ensureDraftID: 99,
	}
	svc := application.NewDraftService(repo)
	meta, err := svc.GetMeta(context.Background(), principal(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if meta.DraftVersionID == nil || *meta.DraftVersionID != 99 {
		t.Fatalf("expected draft 99, got %+v", meta.DraftVersionID)
	}
}

func TestSaveDraft_NotDraftError(t *testing.T) {
	repo := &fakeRepo{saveErr: postgres.ErrNotDraftVersion}
	svc := application.NewDraftService(repo)
	err := svc.SaveDraft(context.Background(), principal(), 1, 2, cfgdomain.Configuration{})
	if err != application.ErrNotDraftVersion {
		t.Fatalf("expected not draft error, got %v", err)
	}
}

func intPtr(n int) *int { return &n }
