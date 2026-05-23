package application

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/port"
)

var ErrMainAppURLMissing = fmt.Errorf("main application download url not configured")

// ProvisioningBuilder builds Android Device Owner QR JSON matching Java QRCodeResource.
type ProvisioningBuilder struct {
	BaseURL         string
	BaseURLForQR    string
	FilesDirectory  string
	ServerProject   string
	SingleCustomer  bool
}

func (b *ProvisioningBuilder) Build(cfg *port.QRConfig, q domain.QRQuery) (string, error) {
	apkURL := ResolveMainAppDownloadURL(cfg, b.BaseURL)
	apkURL = rewriteLoopback(apkURL, b.BaseURL)
	apkURL = strings.ReplaceAll(apkURL, " ", "%20")
	if strings.TrimSpace(apkURL) == "" {
		return "", ErrMainAppURLMissing
	}
	hashURL := ResolveVersionDownloadURL(cfg, b.BaseURL)
	if hashURL == "" {
		hashURL = apkURL
	}
	checksum, err := ApkChecksum(hashURL, b.BaseURL, b.FilesDirectory, cfg.ApkHash)
	if err != nil {
		return "", err
	}
	pkg := cfg.MainAppPkg
	if pkg == "" {
		pkg = "com.hmdm.launcher"
	}
	receiver := strings.TrimSpace(cfg.EventReceivingComponent)
	if receiver == "" {
		receiver = "com.hmdm.launcher.AdminReceiver"
	}
	component := pkg + "/" + receiver

	inner, err := b.adminExtrasBundle(cfg, q)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf("\"android.app.extra.PROVISIONING_DEVICE_ADMIN_COMPONENT_NAME\":%s,\n", jsonQuote(component)))
	sb.WriteString(fmt.Sprintf("\"android.app.extra.PROVISIONING_DEVICE_ADMIN_PACKAGE_DOWNLOAD_LOCATION\":%s,\n", jsonQuote(apkURL)))
	sb.WriteString(fmt.Sprintf("\"android.app.extra.PROVISIONING_DEVICE_ADMIN_PACKAGE_CHECKSUM\":%s,\n", jsonQuote(checksum)))

	if ssid := strings.TrimSpace(cfg.WifiSSID); ssid != "" {
		sec := strings.TrimSpace(cfg.WifiSecurityType)
		if sec == "" {
			sec = "WPA"
		}
		sb.WriteString(fmt.Sprintf("\"android.app.extra.PROVISIONING_WIFI_SSID\":%s,\n", jsonQuote(ssid)))
		sb.WriteString(fmt.Sprintf("\"android.app.extra.PROVISIONING_WIFI_SECURITY_TYPE\":%s,\n", jsonQuote(sec)))
	}
	if pw := strings.TrimSpace(cfg.WifiPassword); pw != "" {
		sb.WriteString(fmt.Sprintf("\"android.app.extra.PROVISIONING_WIFI_PASSWORD\":%s,\n", jsonQuote(pw)))
	}
	if cfg.MobileEnrollment {
		sb.WriteString("\"android.app.extra.PROVISIONING_USE_MOBILE_DATA\":true,\n")
	}
	if extra := strings.TrimSpace(cfg.QRParameters); extra != "" {
		if !strings.HasSuffix(extra, ",") {
			extra += ","
		}
		sb.WriteString(extra + "\n")
	}
	sb.WriteString("\"android.app.extra.PROVISIONING_LEAVE_ALL_SYSTEM_APPS_ENABLED\":true,\n")
	if !cfg.EncryptDevice {
		sb.WriteString("\"android.app.extra.PROVISIONING_SKIP_ENCRYPTION\":true,\n")
	}
	sb.WriteString("\"android.app.extra.PROVISIONING_ADMIN_EXTRAS_BUNDLE\": ")
	sb.WriteString(inner)
	sb.WriteString("\n}\n")
	return sb.String(), nil
}

func (b *ProvisioningBuilder) adminExtrasBundle(cfg *port.QRConfig, q domain.QRQuery) (string, error) {
	m := make(map[string]string)
	if id := strings.TrimSpace(q.DeviceID); id != "" {
		m["com.hmdm.DEVICE_ID"] = id
	}
	if createOnDemand(q.CreateOnDemand) {
		m["com.hmdm.CONFIG"] = cfg.QRCodeKey
		if !b.SingleCustomer && cfg.CustomerName != "" {
			m["com.hmdm.CUSTOMER"] = cfg.CustomerName
		}
	}
	if len(q.Groups) > 0 {
		m["com.hmdm.GROUP"] = strings.Join(q.Groups, ",")
	}
	if mode := strings.TrimSpace(q.UseID); mode != "" {
		m["com.hmdm.DEVICE_ID_USE"] = mode
	}
	host := b.BaseURLForQR
	if host == "" {
		host = strings.TrimRight(b.BaseURL, "/")
	}
	m["com.hmdm.BASE_URL"] = host
	m["com.hmdm.SERVER_PROJECT"] = strings.TrimPrefix(strings.TrimSpace(b.ServerProject), "/")

	inner, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	s := string(inner)
	if extra := strings.TrimSpace(cfg.AdminExtras); extra != "" {
		extra = strings.TrimPrefix(extra, "{")
		extra = strings.TrimSuffix(extra, "}")
		extra = strings.TrimSpace(extra)
		if extra != "" {
			s = strings.TrimSuffix(s, "}") + ",\n" + extra + "\n}"
		}
	}
	return s, nil
}

func createOnDemand(v string) bool {
	v = strings.TrimSpace(v)
	return v == "1" || strings.EqualFold(v, "true")
}

func jsonQuote(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func baseURLHost(base string) string {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	u, err := url.Parse(base)
	if err != nil || u.Host == "" {
		return base
	}
	if u.Port() != "" {
		return u.Scheme + "://" + u.Hostname() + ":" + u.Port()
	}
	return u.Scheme + "://" + u.Hostname()
}
