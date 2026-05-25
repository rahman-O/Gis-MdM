package application

import (
	"strings"
	"testing"
)

func strPtr(s string) *string { return &s }

func TestValidateProvisioningFields_WifiSsid(t *testing.T) {
	// 32 chars OK
	s32 := strings.Repeat("a", 32)
	if err := ValidateProvisioningFields(&s32, nil, nil, nil, nil); err != nil {
		t.Fatalf("32 chars should pass: %v", err)
	}
	// 33 chars fails
	s33 := strings.Repeat("a", 33)
	if err := ValidateProvisioningFields(&s33, nil, nil, nil, nil); err != ErrWifiSsidTooLong {
		t.Fatalf("33 chars should fail with ErrWifiSsidTooLong, got: %v", err)
	}
}

func TestValidateProvisioningFields_WifiPassword(t *testing.T) {
	s63 := strings.Repeat("x", 63)
	if err := ValidateProvisioningFields(nil, &s63, nil, nil, nil); err != nil {
		t.Fatalf("63 chars should pass: %v", err)
	}
	s64 := strings.Repeat("x", 64)
	if err := ValidateProvisioningFields(nil, &s64, nil, nil, nil); err != ErrWifiPasswordTooLong {
		t.Fatalf("64 chars should fail with ErrWifiPasswordTooLong, got: %v", err)
	}
}

func TestValidateProvisioningFields_SecurityType(t *testing.T) {
	valid := []string{"", "WPA", "WPA2", "WPA3", "WEP", "NONE"}
	for _, v := range valid {
		if err := ValidateProvisioningFields(nil, nil, strPtr(v), nil, nil); err != nil {
			t.Fatalf("security type %q should pass: %v", v, err)
		}
	}
	invalid := "INVALID"
	if err := ValidateProvisioningFields(nil, nil, &invalid, nil, nil); err != ErrInvalidSecurityType {
		t.Fatalf("invalid security type should fail, got: %v", err)
	}
}

func TestValidateProvisioningFields_QrParamsJSON(t *testing.T) {
	valid := `{"key":"value"}`
	if err := ValidateProvisioningFields(nil, nil, nil, &valid, nil); err != nil {
		t.Fatalf("valid JSON should pass: %v", err)
	}
	invalid := `{not json`
	if err := ValidateProvisioningFields(nil, nil, nil, &invalid, nil); err != ErrQrParamsInvalidJSON {
		t.Fatalf("invalid JSON should fail with ErrQrParamsInvalidJSON, got: %v", err)
	}
	// Empty string passes
	empty := ""
	if err := ValidateProvisioningFields(nil, nil, nil, &empty, nil); err != nil {
		t.Fatalf("empty string should pass: %v", err)
	}
}

func TestValidateProvisioningFields_AdminExtrasJSON(t *testing.T) {
	valid := `{"extra":"data"}`
	if err := ValidateProvisioningFields(nil, nil, nil, nil, &valid); err != nil {
		t.Fatalf("valid JSON should pass: %v", err)
	}
	invalid := `broken`
	if err := ValidateProvisioningFields(nil, nil, nil, nil, &invalid); err != ErrAdminExtrasInvalidJSON {
		t.Fatalf("invalid JSON should fail with ErrAdminExtrasInvalidJSON, got: %v", err)
	}
}

func TestValidateProvisioningFields_AllNil(t *testing.T) {
	if err := ValidateProvisioningFields(nil, nil, nil, nil, nil); err != nil {
		t.Fatalf("all nil should pass: %v", err)
	}
}

func TestValidateProvisioningFields_AllEmpty(t *testing.T) {
	e := ""
	if err := ValidateProvisioningFields(&e, &e, &e, &e, &e); err != nil {
		t.Fatalf("all empty should pass: %v", err)
	}
}
