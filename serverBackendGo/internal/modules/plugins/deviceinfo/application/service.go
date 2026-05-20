package application

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/deviceinfo/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/deviceinfo/domain"
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
	if p == nil || !p.CanPluginDeviceinfoAccess() {
		return domain.Settings{}, ErrPermissionDenied
	}
	return s.repo.GetSettings(ctx, int64(p.CustomerID))
}

func (s *Service) SaveSettings(ctx context.Context, p *platformauth.Principal, in domain.Settings) error {
	if p == nil || !p.CanPluginDeviceinfoAccess() {
		return ErrPermissionDenied
	}
	in.CustomerID = int64(p.CustomerID)
	return s.repo.SaveSettings(ctx, in)
}

func (s *Service) SavePublicDynamic(ctx context.Context, deviceNumber string, items []domain.DynamicInfo) error {
	id, cid, err := s.repo.DeviceByNumber(ctx, deviceNumber)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrDeviceNotFound
	}
	if err != nil {
		return err
	}
	return s.repo.SaveDynamic(ctx, id, cid, items)
}

func (s *Service) GetDeviceDetail(ctx context.Context, p *platformauth.Principal, deviceNumber string) (domain.DeviceDetail, error) {
	if p == nil || !p.CanPluginDeviceinfoAccess() {
		return domain.DeviceDetail{}, ErrPermissionDenied
	}
	id, cid, err := s.repo.DeviceByNumber(ctx, deviceNumber)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.DeviceDetail{}, ErrDeviceNotFound
	}
	if err != nil {
		return domain.DeviceDetail{}, err
	}
	if cid != int64(p.CustomerID) {
		return domain.DeviceDetail{}, ErrPermissionDenied
	}
	recs, err := s.repo.ListRecords(ctx, id, 50)
	if err != nil {
		return domain.DeviceDetail{}, err
	}
	return domain.DeviceDetail{DeviceNumber: deviceNumber, Records: recs}, nil
}

func (s *Service) SearchDynamic(ctx context.Context, p *platformauth.Principal, f domain.DynamicSearchFilter) (domain.PaginatedDynamic, error) {
	if p == nil || !p.CanPluginDeviceinfoAccess() {
		return domain.PaginatedDynamic{}, ErrPermissionDenied
	}
	recs, err := s.repo.ListRecords(ctx, f.DeviceID, f.PageSize)
	if err != nil {
		return domain.PaginatedDynamic{}, err
	}
	return domain.PaginatedDynamic{Items: recs, Total: int64(len(recs))}, nil
}
