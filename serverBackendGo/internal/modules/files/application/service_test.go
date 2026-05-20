package application

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/files/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

func TestSearchPermissionDenied(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, "http://localhost:8080", nil)
	_, err := svc.Search(context.Background(), &platformauth.Principal{Permissions: []string{}}, "")
	if err != ErrPermissionDenied {
		t.Fatalf("got %v", err)
	}
}

func TestRemoveFileUsed(t *testing.T) {
	files := &stubFiles{usedCfg: true}
	svc := NewService(files, &stubCustomer{}, nil, nil, "http://localhost:8080", nil)
	p := &platformauth.Principal{CustomerID: 1, Permissions: []string{"edit_files"}}
	id := 1
	err := svc.Remove(context.Background(), p, id, "a.txt", false)
	if err != ErrFileUsed {
		t.Fatalf("got %v", err)
	}
}

type stubFiles struct {
	usedCfg bool
}

func (stubFiles) List(context.Context, int, string) ([]domain.UploadedFile, error) { return nil, nil }
func (stubFiles) GetByID(context.Context, int, int) (*domain.UploadedFile, error) { return nil, nil }
func (stubFiles) Insert(context.Context, *domain.UploadedFile) error { return nil }
func (stubFiles) Update(context.Context, *domain.UploadedFile) error { return nil }
func (stubFiles) Delete(context.Context, int) error { return nil }
func (s stubFiles) IsUsedByConfiguration(context.Context, int) (bool, error) { return s.usedCfg, nil }
func (stubFiles) IsUsedByIcon(context.Context, int) (bool, error) { return false, nil }
func (stubFiles) UsingConfigurationNames(context.Context, int, int) ([]string, error) { return nil, nil }
func (stubFiles) UsingIconNames(context.Context, int, int) ([]string, error) { return nil, nil }
func (stubFiles) GetFileConfigurations(context.Context, int, int, int) ([]domain.FileConfigurationLink, error) {
	return nil, nil
}
func (stubFiles) DeleteConfigurationFile(context.Context, int) error { return nil }
func (stubFiles) InsertConfigurationFile(context.Context, int, int, string) error { return nil }
func (stubFiles) CountByPath(context.Context, int, *int, string) (int64, error) { return 0, nil }

type stubCustomer struct{}

func (stubCustomer) GetMeta(context.Context, int) (*domain.CustomerMeta, error) {
	return &domain.CustomerMeta{FilesDir: "c1"}, nil
}
func (stubCustomer) CountCustomers(context.Context) (int, error) { return 1, nil }
