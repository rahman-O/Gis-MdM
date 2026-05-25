package application

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrWifiSsidTooLong       = errors.New("error.enrollment_route.wifi_ssid_too_long")
	ErrWifiPasswordTooLong   = errors.New("error.enrollment_route.wifi_password_too_long")
	ErrInvalidSecurityType   = errors.New("error.enrollment_route.invalid_security_type")
	ErrQrParamsInvalidJSON   = errors.New("error.enrollment_route.qr_parameters_invalid_json")
	ErrQrParamsTooLong       = errors.New("error.enrollment_route.qr_parameters_too_long")
	ErrAdminExtrasInvalidJSON = errors.New("error.enrollment_route.admin_extras_invalid_json")
	ErrAdminExtrasTooLong    = errors.New("error.enrollment_route.admin_extras_too_long")
)

// ValidateProvisioningFields checks provisioning field constraints.
// All fields are optional (nil = skip validation for that field).
func ValidateProvisioningFields(ssid, password, secType, qrParams, adminExtras *string) error {
	if ssid != nil && len(*ssid) > 32 {
		return ErrWifiSsidTooLong
	}
	if password != nil && len(*password) > 63 {
		return ErrWifiPasswordTooLong
	}
	if secType != nil && !isValidSecurityType(*secType) {
		return ErrInvalidSecurityType
	}
	if qrParams != nil && strings.TrimSpace(*qrParams) != "" {
		if len(*qrParams) > 65535 {
			return ErrQrParamsTooLong
		}
		if !json.Valid([]byte(*qrParams)) {
			return ErrQrParamsInvalidJSON
		}
	}
	if adminExtras != nil && strings.TrimSpace(*adminExtras) != "" {
		if len(*adminExtras) > 65535 {
			return ErrAdminExtrasTooLong
		}
		if !json.Valid([]byte(*adminExtras)) {
			return ErrAdminExtrasInvalidJSON
		}
	}
	return nil
}

func isValidSecurityType(s string) bool {
	switch strings.TrimSpace(s) {
	case "", "WPA", "WPA2", "WPA3", "WEP", "NONE":
		return true
	}
	return false
}
