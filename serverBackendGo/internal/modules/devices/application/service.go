package application

import (
	"context"
	"errors"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/devices/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/devices/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// Service implements device use cases.
type Service struct {
	repo port.DeviceRepository
	push port.PushNotifier
}

func NewService(repo port.DeviceRepository, push port.PushNotifier) *Service {
	if push == nil {
		push = port.NoopPush{}
	}
	return &Service{repo: repo, push: push}
}

var (
	ErrPermissionDenied = errors.New("error.permission.denied")
	ErrDeviceExists     = errors.New("error.duplicate.device")
	ErrDeviceNotFound   = errors.New("error.notfound.object")
	ErrDeviceLimit      = errors.New("error.device.limit")
)

func (s *Service) scope(ctx context.Context, p *platformauth.Principal) (*port.UserScope, error) {
	if p == nil {
		return nil, ErrPermissionDenied
	}
	return s.repo.LoadUserScope(ctx, p.ID)
}

func (s *Service) Search(ctx context.Context, p *platformauth.Principal, req domain.SearchRequest) (*domain.DeviceListView, error) {
	scope, err := s.scope(ctx, p)
	if err != nil {
		return nil, err
	}
	req.Prepare()
	items, err := s.repo.Search(ctx, *scope, req)
	if err != nil {
		return nil, err
	}
	total, err := s.repo.Count(ctx, *scope, req)
	if err != nil {
		return nil, err
	}
	configs, err := s.repo.ListConfigurations(ctx, scope.CustomerID)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.DeviceView{}
	}
	return &domain.DeviceListView{
		Configurations: configs,
		Devices: domain.DevicePage{
			Items:           items,
			TotalItemsCount: total,
		},
	}, nil
}

func (s *Service) GetByNumber(ctx context.Context, p *platformauth.Principal, number string) (*domain.DeviceView, error) {
	scope, err := s.scope(ctx, p)
	if err != nil {
		return nil, err
	}
	d, err := s.repo.GetByNumber(ctx, *scope, number)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, ErrDeviceNotFound
	}
	return d, nil
}

func (s *Service) Save(ctx context.Context, p *platformauth.Principal, d domain.SaveDevice) error {
	if !p.CanEditDevices() {
		return ErrPermissionDenied
	}
	scope, err := s.scope(ctx, p)
	if err != nil {
		return err
	}
	if len(d.IDs) > 0 && d.ConfigurationID != nil {
		return s.repo.UpdateConfigurationBulk(ctx, scope.CustomerID, d.IDs, *d.ConfigurationID)
	}
	number := strings.TrimSpace(ptrStr(d.Number))
	if number == "" {
		return ErrPermissionDenied
	}
	exclude := 0
	if d.ID != nil {
		exclude = *d.ID
	}
	exists, err := s.repo.ExistsNumber(ctx, scope.CustomerID, number, exclude)
	if err != nil {
		return err
	}
	if exists {
		return ErrDeviceExists
	}
	if d.ID == nil || *d.ID == 0 {
		limit, err := s.repo.DeviceLimit(ctx, scope.CustomerID)
		if err != nil {
			return err
		}
		if limit > 0 {
			count, err := s.repo.CountDevices(ctx, scope.CustomerID)
			if err != nil {
				return err
			}
			if count >= int64(limit) {
				return ErrDeviceLimit
			}
		}
		if d.ConfigurationID == nil {
			return ErrPermissionDenied
		}
		_, err = s.repo.Insert(ctx, scope.CustomerID, d)
		return err
	}
	return s.repo.Update(ctx, scope.CustomerID, d)
}

func (s *Service) Delete(ctx context.Context, p *platformauth.Principal, id int) error {
	if !p.CanEditDevices() {
		return ErrPermissionDenied
	}
	scope, err := s.scope(ctx, p)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, scope.CustomerID, id)
}

func (s *Service) DeleteBulk(ctx context.Context, p *platformauth.Principal, req domain.BulkDeleteRequest) error {
	if !p.CanEditDevices() {
		return ErrPermissionDenied
	}
	scope, err := s.scope(ctx, p)
	if err != nil {
		return err
	}
	return s.repo.DeleteBulk(ctx, scope.CustomerID, req.IDs)
}

func (s *Service) GroupBulk(ctx context.Context, p *platformauth.Principal, req domain.GroupBulkRequest) error {
	if !p.CanEditDevices() {
		return ErrPermissionDenied
	}
	scope, err := s.scope(ctx, p)
	if err != nil {
		return err
	}
	return s.repo.UpdateGroupBulk(ctx, scope.CustomerID, req)
}

func (s *Service) Autocomplete(ctx context.Context, p *platformauth.Principal, filter string) ([]domain.LookupItem, error) {
	scope, err := s.scope(ctx, p)
	if err != nil {
		return nil, err
	}
	return s.repo.Autocomplete(ctx, *scope, filter, 10)
}

func (s *Service) UpdateDescription(ctx context.Context, p *platformauth.Principal, id int, description string) error {
	if !p.CanEditDeviceDescription() {
		return ErrPermissionDenied
	}
	scope, err := s.scope(ctx, p)
	if err != nil {
		return err
	}
	return s.repo.UpdateDescription(ctx, scope.CustomerID, id, description)
}

func (s *Service) GetAppSettings(ctx context.Context, p *platformauth.Principal, deviceID int) ([]domain.AppSetting, error) {
	if p == nil {
		return nil, ErrPermissionDenied
	}
	return s.repo.ListAppSettings(ctx, deviceID)
}

func (s *Service) SaveAppSettings(ctx context.Context, p *platformauth.Principal, deviceID int, settings []domain.AppSetting) error {
	if !p.CanEditDevices() {
		return ErrPermissionDenied
	}
	return s.repo.SaveAppSettings(ctx, deviceID, settings)
}

func (s *Service) NotifyAppSettings(ctx context.Context, deviceID int) error {
	return s.push.NotifyAppSettings(ctx, deviceID)
}

func (s *Service) MoveTree(ctx context.Context, p *platformauth.Principal, deviceID int, treeNodeID int) error {
	if !p.CanEditDevices() {
		return ErrPermissionDenied
	}
	scope, err := s.scope(ctx, p)
	if err != nil {
		return err
	}
	if err := s.repo.MoveTreeNode(ctx, scope.CustomerID, deviceID, treeNodeID); err != nil {
		return ErrDeviceNotFound
	}
	return nil
}

func ptrStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
