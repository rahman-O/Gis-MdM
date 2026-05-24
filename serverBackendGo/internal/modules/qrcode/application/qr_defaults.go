package application

import (
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/port"
)

// applyEnrollmentRouteQRDefaults enforces create=1 and default device id mode for public enrollment QR.
func applyEnrollmentRouteQRDefaults(q domain.QRQuery, cfg *port.QRConfig) domain.QRQuery {
	if strings.TrimSpace(q.DeviceID) != "" {
		return q
	}
	if !createOnDemand(q.CreateOnDemand) {
		q.CreateOnDemand = "1"
	}
	if strings.TrimSpace(q.UseID) == "" {
		mode := strings.TrimSpace(cfg.DefaultDeviceIDMode)
		if mode == "" {
			mode = "imei"
		}
		q.UseID = mode
	}
	return q
}
