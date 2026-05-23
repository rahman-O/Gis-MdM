package domain

import (
	"encoding/json"
	"testing"
)

func TestConfigurationPolicyRoundTrip(t *testing.T) {
	cfg := Configuration{
		Name: ptr("Test"),
		Policy: map[string]any{
			"kioskMode": true,
			"wifi":      false,
		},
	}
	raw, err := cfg.BuildSettingsJSON()
	if err != nil {
		t.Fatal(err)
	}
	var loaded Configuration
	loaded.SetPolicyFromJSON(raw)
	if loaded.Policy["kioskMode"] != true {
		t.Fatalf("kioskMode: %v", loaded.Policy["kioskMode"])
	}
	m := ConfigurationResponseMap(&cfg)
	if m["kioskMode"] != true {
		t.Fatalf("response map: %v", m["kioskMode"])
	}
}

func TestParseConfigurationBody(t *testing.T) {
	body := []byte(`{"name":"X","kioskMode":true,"applications":[]}`)
	cfg, err := ParseConfigurationBody(body)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Policy["kioskMode"] != true {
		t.Fatalf("policy: %v", cfg.Policy)
	}
	b, _ := json.Marshal(cfg)
	_ = b
}

func ptr(s string) *string { return &s }
