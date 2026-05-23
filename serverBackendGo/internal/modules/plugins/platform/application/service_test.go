package application

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/config"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/shared/status"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

func TestSaveDisabledPermission(t *testing.T) {
	svc := NewService(nil, status.NewCache())
	p := &platformauth.Principal{CustomerID: 1}
	err := svc.SaveDisabled(context.Background(), p, []int64{1})
	if err != ErrPermissionDenied {
		t.Fatalf("want permission denied, got %v", err)
	}
}

func TestEnabledPluginsFilter(t *testing.T) {
	cfg := config.Config{EnabledPlugins: []string{"audit", "push"}}
	if !cfg.IsPluginEnabled("audit") {
		t.Fatal("audit should be enabled")
	}
	if cfg.IsPluginEnabled("xtra") {
		t.Fatal("xtra should be disabled")
	}
}
