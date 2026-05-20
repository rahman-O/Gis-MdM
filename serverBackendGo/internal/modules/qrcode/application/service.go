package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/skip2/go-qrcode"
	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/port"
)

var ErrNotFound = errors.New("configuration not found")

type Service struct {
	repo    port.ConfigByKey
	baseURL string
}

func NewService(repo port.ConfigByKey, baseURL string) *Service {
	return &Service{repo: repo, baseURL: strings.TrimRight(baseURL, "/")}
}

func (s *Service) extrasBundle(cfg *port.QRConfig, q domain.QRQuery) string {
	_ = q
	url := cfg.MainAppURL
	if cfg.LauncherURL != "" {
		url = cfg.LauncherURL
	}
	url = rewriteLoopback(url, s.baseURL)
	pkg := cfg.MainAppPkg
	if pkg == "" {
		pkg = "com.hmdm.launcher"
	}
	extra := strings.TrimSpace(cfg.AdminExtras)
	if extra != "" {
		extra = ",\n" + extra
	}
	return fmt.Sprintf(`{
  "android.app.extra.PROVISIONING_DEVICE_ADMIN_COMPONENT_NAME":"%s/%s.AdminReceiver",
  "android.app.extra.PROVISIONING_DEVICE_ADMIN_PACKAGE_DOWNLOAD_LOCATION":"%s"%s
}`, pkg, pkg, url, extra)
}

func rewriteLoopback(url, base string) string {
	if url == "" {
		return url
	}
	if strings.Contains(url, "127.0.0.1") || strings.Contains(url, "localhost") {
		host := strings.TrimPrefix(strings.TrimPrefix(base, "https://"), "http://")
		if i := strings.Index(host, "/"); i >= 0 {
			host = host[:i]
		}
		return strings.ReplaceAll(strings.ReplaceAll(url, "127.0.0.1", host), "localhost", host)
	}
	return url
}

func (s *Service) JSON(ctx context.Context, key string, q domain.QRQuery) (string, error) {
	cfg, err := s.repo.ConfigurationByQRKey(ctx, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	return s.extrasBundle(cfg, q), nil
}

func (s *Service) PNG(ctx context.Context, key string, q domain.QRQuery) ([]byte, error) {
	body, err := s.JSON(ctx, key, q)
	if err != nil {
		return nil, err
	}
	size := q.Size
	if size <= 0 {
		size = 250
	}
	return qrcode.Encode(body, qrcode.Medium, size)
}
