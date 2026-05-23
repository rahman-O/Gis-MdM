package application

import (
	"encoding/json"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/sync/domain"
)

func TestApplyConfigurationPolicy_kioskAndRestrictions(t *testing.T) {
	settings := []byte(`{"kioskMode":true,"restrictions":"no_usb","gps":false,"pushOptions":"mqtt"}`)
	resp := &domain.SyncResponse{}
	ApplyConfigurationPolicy(resp, settings, nil)
	if !resp.KioskMode {
		t.Fatal("expected kioskMode true")
	}
	if resp.Restrictions == nil || *resp.Restrictions != "no_usb" {
		t.Fatalf("restrictions: %v", resp.Restrictions)
	}
	if resp.GPS == nil || *resp.GPS {
		t.Fatalf("gps: %v", resp.GPS)
	}
	if resp.PushOptions == nil || *resp.PushOptions != "mqtt" {
		t.Fatalf("pushOptions: %v", resp.PushOptions)
	}
}

func TestApplyConfigurationPolicy_backgroundImage(t *testing.T) {
	url := "https://example.com/bg.png"
	resp := &domain.SyncResponse{}
	ApplyConfigurationPolicy(resp, nil, &url)
	if resp.BackgroundImageURL == nil || *resp.BackgroundImageURL != url {
		t.Fatalf("bg: %v", resp.BackgroundImageURL)
	}
}

func TestApplyConfigurationPolicy_iconSizeEnumToInt(t *testing.T) {
	settings := []byte(`{"iconSize":"SMALL"}`)
	resp := &domain.SyncResponse{}
	ApplyConfigurationPolicy(resp, settings, nil)
	if resp.IconSize == nil || *resp.IconSize != 100 {
		t.Fatalf("iconSize: %v", resp.IconSize)
	}
}

func TestApplyConfigurationPolicy_invalidJSON(t *testing.T) {
	resp := &domain.SyncResponse{}
	ApplyConfigurationPolicy(resp, []byte(`{`), nil)
	var _ = json.Valid([]byte(`{}`))
}
