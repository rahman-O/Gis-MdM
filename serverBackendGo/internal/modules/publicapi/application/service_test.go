package application

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/publicapi/domain"
	sharedcrypto "github.com/gis-mdm/server-backend-go/internal/shared/crypto"
)

func TestUploadInvalidHash(t *testing.T) {
	svc := NewService(&stubDevice{}, nil, "http://localhost:8080", RebrandingConfig{HashSecret: "secret"})
	err := svc.UploadApplication(context.Background(), `{"deviceId":"d1","hash":"bad","name":"A","pkg":"com.x","version":"1"}`, "", nil)
	if err != ErrInvalidHash {
		t.Fatalf("got %v", err)
	}
}

func TestUploadValidHashMissingDevice(t *testing.T) {
	hash := sharedcrypto.DeviceUploadHash("d1", "secret")
	body := `{"deviceId":"d1","hash":"` + hash + `","name":"A","pkg":"com.x","version":"1"}`
	svc := NewService(&stubDevice{}, nil, "http://localhost:8080", RebrandingConfig{HashSecret: "secret"})
	err := svc.UploadApplication(context.Background(), body, "", nil)
	if err != ErrDeviceNotFound {
		t.Fatalf("got %v", err)
	}
}

type stubDevice struct{}

func (stubDevice) FindDeviceByNumber(context.Context, string) (*domain.DeviceRef, error) {
	return nil, nil
}
func (stubDevice) CustomerFilesDir(context.Context, int) (string, error) { return "c1", nil }
func (stubDevice) HasDuplicateApp(context.Context, int, string, string) (bool, error) {
	return false, nil
}
func (stubDevice) InsertApplication(context.Context, int, string, string, string, string, domain.UploadAppRequest) error {
	return nil
}
