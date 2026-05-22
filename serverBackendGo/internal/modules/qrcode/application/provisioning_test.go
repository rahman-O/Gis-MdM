package application

import (
	"strings"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/port"
)

func TestProvisioningBuilder_containsRequiredKeys(t *testing.T) {
	b := ProvisioningBuilder{
		BaseURL:        "http://192.168.1.10:8080",
		BaseURLForQR:   "http://192.168.1.10:8080",
		FilesDirectory: t.TempDir(),
		SingleCustomer: true,
	}
	cfg := &port.QRConfig{
		QRCodeKey:   "default-qr",
		MainAppPkg:  "com.hmdm.launcher",
		MainAppURL:  "http://192.168.1.10:8080/files/customer-1/apk/launcher.apk",
		ApkHash:     "dGVzdGhhc2g=",
	}
	body, err := b.Build(cfg, domain.QRQuery{DeviceID: "dev-1", CreateOnDemand: "1"})
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{
		"PROVISIONING_DEVICE_ADMIN_PACKAGE_CHECKSUM",
		"PROVISIONING_ADMIN_EXTRAS_BUNDLE",
		"com.hmdm.BASE_URL",
		"com.hmdm.CONFIG",
		"com.hmdm.DEVICE_ID",
	} {
		if !strings.Contains(body, key) {
			t.Fatalf("missing %q in provisioning json:\n%s", key, body)
		}
	}
}
