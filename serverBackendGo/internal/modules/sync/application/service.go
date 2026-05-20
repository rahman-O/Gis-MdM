package application

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/sync/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/sync/port"
)

var (
	ErrDeviceNotFound    = errors.New("error.notfound.device")
	ErrDeviceExists      = errors.New("error.duplicate.device")
	ErrPermissionDenied  = errors.New("error.permission.denied")
	ErrMultiTenantCreate = errors.New("error.permission.denied")
)

type Config struct {
	BaseURL            string
	FilesDirectory     string
	HashSecret         string
	SecureEnrollment   bool
	PreventDuplicate   bool
	MobileAppName      string
	VendorName         string
	DefaultCustomerID  int64
}

type Service struct {
	repo port.SyncRepository
	cfg  Config
}

func NewService(repo port.SyncRepository, cfg Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

func (s *Service) checkSig(signature, deviceID string) bool {
	if !s.cfg.SecureEnrollment {
		return true
	}
	return checkRequestSignature(signature, s.cfg.HashSecret+deviceID)
}

func (s *Service) GetConfiguration(ctx context.Context, deviceID, signature, cpuArch string) (*domain.SyncResponse, error) {
	if !s.checkSig(signature, deviceID) {
		return nil, ErrPermissionDenied
	}
	dev, migration, err := s.resolveDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if migration && dev.OldNumber != nil {
		_ = s.repo.CompleteMigration(ctx, dev.ID)
		dev.OldNumber = nil
	}
	if s.cfg.PreventDuplicate && dev.LastUpdate > 0 {
		return nil, ErrDeviceExists
	}
	_ = s.repo.TouchLastUpdate(ctx, dev.ID)
	return s.repo.BuildSyncResponse(ctx, *dev, s.cfg.BaseURL, s.cfg.FilesDirectory, cpuArch, s.cfg.MobileAppName, s.cfg.VendorName)
}

func (s *Service) EnrollConfiguration(ctx context.Context, deviceID string, opts domain.DeviceCreateOptions, signature, cpuArch string) (*domain.SyncResponse, error) {
	if !s.checkSig(signature, deviceID) {
		return nil, ErrPermissionDenied
	}
	dev, err := s.repo.FindByNumber(ctx, deviceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			created, cerr := s.repo.CreateOnDemand(ctx, deviceID, opts, s.cfg.DefaultCustomerID)
			if cerr != nil {
				return nil, cerr
			}
			dev = created
		} else {
			return nil, err
		}
	}
	if s.cfg.PreventDuplicate && dev.LastUpdate > 0 {
		return nil, ErrDeviceExists
	}
	_ = s.repo.TouchLastUpdate(ctx, dev.ID)
	return s.repo.BuildSyncResponse(ctx, *dev, s.cfg.BaseURL, s.cfg.FilesDirectory, cpuArch, s.cfg.MobileAppName, s.cfg.VendorName)
}

func (s *Service) UpdateInfo(ctx context.Context, info domain.DeviceInfo) error {
	dev, err := s.repo.FindByNumber(ctx, info.DeviceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			n, _ := s.repo.CountCustomers(ctx)
			if n > 1 {
				return ErrMultiTenantCreate
			}
			dev, err = s.repo.CreateOnDemand(ctx, info.DeviceID, domain.DeviceCreateOptions{}, s.cfg.DefaultCustomerID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	b, _ := json.Marshal(info)
	_ = s.repo.UpdateInfo(ctx, dev.ID, string(b), "")
	if info.Custom1 != nil || info.Custom2 != nil || info.Custom3 != nil {
		_ = s.repo.UpdateCustomProps(ctx, dev.ID, info.Custom1, info.Custom2, info.Custom3)
	}
	return nil
}

func (s *Service) SaveApplicationSettings(ctx context.Context, deviceNumber string, settings []domain.SyncApplicationSetting) error {
	dev, err := s.repo.FindByNumber(ctx, deviceNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDeviceNotFound
		}
		return err
	}
	return s.repo.SaveApplicationSettings(ctx, dev.ID, settings)
}

func (s *Service) resolveDevice(ctx context.Context, number string) (*domain.DeviceRecord, bool, error) {
	dev, err := s.repo.FindByNumber(ctx, number)
	if err == nil {
		return dev, false, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, false, err
	}
	dev, err = s.repo.FindByOldNumber(ctx, number)
	if err == nil {
		return dev, true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, ErrDeviceNotFound
	}
	return nil, false, err
}
