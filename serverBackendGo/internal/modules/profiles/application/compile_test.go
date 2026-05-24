package application_test

import (
	"encoding/json"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	syncapp "github.com/gis-mdm/server-backend-go/internal/modules/sync/application"
	syncdomain "github.com/gis-mdm/server-backend-go/internal/modules/sync/domain"
)

func TestApplyArtifactToSyncResponse_mapsPolicyLikeMapper(t *testing.T) {
	settings := []byte(`{"gps":true,"kioskMode":true,"pushOptions":"mqtt"}`)
	artifact := &domain.ProfileArtifact{
		ProfileID: 1, ProfileVersionID: 2, VersionNumber: 1,
		Permissive: true, SettingsJSON: settings,
	}
	resp := &syncdomain.SyncResponse{DeviceID: "dev1", ConfigurationID: 5}
	application.ApplyArtifactToSyncResponse(resp, artifact)

	legacy := &syncdomain.SyncResponse{DeviceID: "dev1", ConfigurationID: 5}
	syncapp.ApplyConfigurationPolicy(legacy, settings, nil)

	if resp.GPS == nil || !*resp.GPS {
		t.Fatal("expected gps from artifact policy")
	}
	if !resp.KioskMode {
		t.Fatal("expected kioskMode true")
	}
	if resp.PushOptions == nil || *resp.PushOptions != "mqtt" {
		t.Fatal("expected pushOptions mqtt")
	}
	if legacy.PushOptions == nil || *legacy.PushOptions != "mqtt" {
		t.Fatal("legacy mapper baseline failed")
	}
}

func TestProfileArtifact_roundTripJSON(t *testing.T) {
	artifact := domain.ProfileArtifact{
		ProfileID: 10, ProfileVersionID: 46, VersionNumber: 4,
		SettingsJSON: json.RawMessage(`{"wifi":false}`),
		Applications: []syncdomain.SyncApplication{{ID: 1, Name: "App", Pkg: "com.test", Type: "app"}},
	}
	raw, err := json.Marshal(artifact)
	if err != nil {
		t.Fatal(err)
	}
	var decoded domain.ProfileArtifact
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.ProfileVersionID != 46 || len(decoded.Applications) != 1 {
		t.Fatalf("unexpected decode: %+v", decoded)
	}
}
