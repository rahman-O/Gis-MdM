package application

import (
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/port"
)

func TestApplyEnrollmentRouteQRDefaults(t *testing.T) {
	cfg := &port.QRConfig{DefaultDeviceIDMode: "imei"}
	q := applyEnrollmentRouteQRDefaults(domain.QRQuery{}, cfg)
	if q.CreateOnDemand != "1" {
		t.Fatalf("expected create=1, got %q", q.CreateOnDemand)
	}
	if q.UseID != "imei" {
		t.Fatalf("expected useId imei, got %q", q.UseID)
	}
}
