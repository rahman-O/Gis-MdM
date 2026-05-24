package application

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	"github.com/skip2/go-qrcode"
	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/qrcode/port"
)

var ErrNotFound = errors.New("configuration not found")

type Service struct {
	repo    port.ConfigByKey
	builder ProvisioningBuilder
	log     *slog.Logger
}

func NewService(repo port.ConfigByKey, baseURL, filesDirectory, serverProject string, log *slog.Logger) *Service {
	if log == nil {
		log = slog.Default()
	}
	return &Service{
		repo: repo,
		builder: ProvisioningBuilder{
			BaseURL:        strings.TrimRight(strings.TrimSpace(baseURL), "/"),
			BaseURLForQR:   baseURLHost(baseURL),
			FilesDirectory: filesDirectory,
			ServerProject:  serverProject,
		},
		log: log,
	}
}

func (s *Service) provisioningJSON(ctx context.Context, key string, q domain.QRQuery) (string, error) {
	cfg, err := s.repo.ConfigurationByQRKey(ctx, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	if n, _ := s.repo.CountCustomers(ctx); n <= 1 {
		s.builder.SingleCustomer = true
	} else {
		s.builder.SingleCustomer = false
	}
	q = applyEnrollmentRouteQRDefaults(q, cfg)
	body, err := s.builder.Build(cfg, q)
	if err != nil {
		s.log.Warn("qr provisioning failed", "key", key, "err", err)
		return "", err
	}
	return body, nil
}

func (s *Service) JSON(ctx context.Context, key string, q domain.QRQuery) (string, error) {
	return s.provisioningJSON(ctx, key, q)
}

func (s *Service) PNG(ctx context.Context, key string, q domain.QRQuery) ([]byte, error) {
	body, err := s.provisioningJSON(ctx, key, q)
	if err != nil {
		return nil, err
	}
	size := q.Size
	if size <= 0 {
		size = 250
	}
	return qrcode.Encode(body, qrcode.Medium, size)
}

func rewriteLoopback(url, base string) string {
	if url == "" {
		return url
	}
	if strings.Contains(url, "127.0.0.1") || strings.Contains(url, "localhost") {
		host := baseURLHost(base)
		host = strings.TrimPrefix(strings.TrimPrefix(host, "https://"), "http://")
		if i := strings.Index(host, "/"); i >= 0 {
			host = host[:i]
		}
		url = strings.ReplaceAll(strings.ReplaceAll(url, "127.0.0.1", host), "localhost", host)
	}
	return url
}
