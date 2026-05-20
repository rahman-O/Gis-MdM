package application

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/devicelog/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/devicelog/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

var (
	ErrPermissionDenied = errors.New("error.permission.denied")
	ErrDeviceNotFound     = errors.New("error.device.notfound")
)

type Service struct {
	repo *postgres.Repository
}

func NewService(repo *postgres.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetSettings(ctx context.Context, p *platformauth.Principal) (domain.Settings, error) {
	if p == nil || !p.CanPluginDevicelogAccess() {
		return domain.Settings{}, ErrPermissionDenied
	}
	return s.repo.GetSettings(ctx, int64(p.CustomerID))
}

func (s *Service) SaveSettings(ctx context.Context, p *platformauth.Principal, in domain.Settings) error {
	if p == nil || !p.CanPluginDevicelogAccess() {
		return ErrPermissionDenied
	}
	in.CustomerID = int64(p.CustomerID)
	return s.repo.SaveSettings(ctx, in)
}

func (s *Service) SaveRule(ctx context.Context, p *platformauth.Principal, rule domain.Rule) (int64, error) {
	if p == nil || !p.CanPluginDevicelogAccess() {
		return 0, ErrPermissionDenied
	}
	settings, err := s.repo.GetSettings(ctx, int64(p.CustomerID))
	if err != nil {
		return 0, err
	}
	rule.SettingID = settings.ID
	return s.repo.UpsertRule(ctx, rule)
}

func (s *Service) DeleteRule(ctx context.Context, p *platformauth.Principal, ruleID int64) error {
	if p == nil || !p.CanPluginDevicelogAccess() {
		return ErrPermissionDenied
	}
	return s.repo.DeleteRule(ctx, int64(p.CustomerID), ruleID)
}

func (s *Service) UploadLogs(ctx context.Context, deviceNumber string, rows []domain.UploadRecord) error {
	deviceID, customerID, err := s.repo.DeviceByNumber(ctx, deviceNumber)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrDeviceNotFound
	}
	if err != nil {
		return err
	}
	appID, err := s.repo.DefaultApplicationID(ctx, customerID)
	if err != nil {
		return err
	}
	return s.repo.InsertLogs(ctx, customerID, deviceID, appID, rows)
}

func (s *Service) SearchLogs(ctx context.Context, p *platformauth.Principal, f domain.LogFilter) (domain.PaginatedLogs, error) {
	if p == nil || !p.CanPluginDevicelogAccess() {
		return domain.PaginatedLogs{}, ErrPermissionDenied
	}
	items, total, err := s.repo.SearchLogs(ctx, int64(p.CustomerID), f)
	if err != nil {
		return domain.PaginatedLogs{}, err
	}
	return domain.PaginatedLogs{Items: items, Total: total}, nil
}
