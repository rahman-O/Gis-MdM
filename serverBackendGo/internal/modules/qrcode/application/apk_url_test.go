package application

import (
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/port"
)

func TestResolveMainAppDownloadURL_priority(t *testing.T) {
	cfg := &port.QRConfig{
		LauncherURL:     "http://lan/launcher.apk",
		MainAppURL:      "http://ignored/version.apk",
		MainAppFilePath: "app/demo.apk",
		FilesDir:        "customer-1",
	}
	if got := ResolveMainAppDownloadURL(cfg, "http://192.168.1.5:8080"); got != cfg.LauncherURL {
		t.Fatalf("launcher override: got %q", got)
	}
	cfg.LauncherURL = ""
	if got := ResolveMainAppDownloadURL(cfg, "http://192.168.1.5:8080"); got != cfg.MainAppURL {
		t.Fatalf("version url: got %q", got)
	}
	cfg.MainAppURL = ""
	cfg.AppLevelURL = "http://app-level.apk"
	if got := ResolveMainAppDownloadURL(cfg, "http://192.168.1.5:8080"); got != cfg.AppLevelURL {
		t.Fatalf("app url: got %q", got)
	}
	cfg.AppLevelURL = ""
	got := ResolveMainAppDownloadURL(cfg, "http://192.168.1.5:8080")
	want := "http://192.168.1.5:8080/files/customer-1/app/demo.apk"
	if got != want {
		t.Fatalf("filepath build: got %q want %q", got, want)
	}
}
