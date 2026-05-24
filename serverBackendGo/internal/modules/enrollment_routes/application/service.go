package application

import (
	"context"
	"errors"
	"strings"

	cfgdomain "github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type Service struct {
	repo port.RouteRepository
}

func NewService(repo port.RouteRepository) *Service {
	return &Service{repo: repo}
}

var (
	ErrPermissionDenied       = errors.New("error.permission.denied")
	ErrRouteNotFound          = errors.New("error.notfound.enrollment_route")
	ErrDuplicateRoute         = errors.New("error.duplicate.enrollment_route")
	ErrPublishedVersionRequired = errors.New("error.enrollment_route.published_version_required")
	ErrTreeNodeRequired       = errors.New("error.enrollment_route.tree_node_required")
	ErrMainAppRequired        = errors.New("error.enrollment_route.main_app_required")
)

func customerID(p *platformauth.Principal) int {
	if p == nil {
		return 0
	}
	return p.CustomerID
}

func (s *Service) requirePerm(p *platformauth.Principal) error {
	if p == nil || !p.CanManageConfigurations() {
		return ErrPermissionDenied
	}
	return nil
}

func (s *Service) List(ctx context.Context, p *platformauth.Principal) ([]domain.Route, error) {
	if err := s.requirePerm(p); err != nil {
		return nil, err
	}
	items, err := s.repo.List(ctx, customerID(p))
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.Route{}
	}
	return items, nil
}

func (s *Service) GetByID(ctx context.Context, p *platformauth.Principal, id int) (*domain.RouteDetail, error) {
	if err := s.requirePerm(p); err != nil {
		return nil, err
	}
	detail, err := s.repo.GetByID(ctx, customerID(p), id)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, ErrRouteNotFound
	}
	return detail, nil
}

func (s *Service) ListPublishedProfileVersions(ctx context.Context, p *platformauth.Principal) ([]domain.PublishedProfileVersion, error) {
	if err := s.requirePerm(p); err != nil {
		return nil, err
	}
	items, err := s.repo.ListPublishedProfileVersions(ctx, customerID(p))
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.PublishedProfileVersion{}
	}
	return items, nil
}

func (s *Service) Create(ctx context.Context, p *platformauth.Principal, req domain.CreateRequest) (*domain.RouteDetail, error) {
	if err := s.requirePerm(p); err != nil {
		return nil, err
	}
	pvID := 0
	if req.ProfileVersionID != nil {
		pvID = *req.ProfileVersionID
	}
	if err := s.validateBinding(ctx, customerID(p), pvID, req.DefaultTreeNodeID, req.MainAppID); err != nil {
		return nil, err
	}
	key := cfgdomain.NewQRCodeKey()
	id, err := s.repo.Create(ctx, customerID(p), req, key)
	if err != nil {
		if errors.Is(err, postgres.ErrDuplicateName) {
			return nil, ErrDuplicateRoute
		}
		return nil, err
	}
	return s.GetByID(ctx, p, id)
}

func (s *Service) Update(ctx context.Context, p *platformauth.Principal, id int, req domain.UpdateRequest) (*domain.RouteDetail, error) {
	if err := s.requirePerm(p); err != nil {
		return nil, err
	}
	cur, err := s.repo.GetByID(ctx, customerID(p), id)
	if err != nil {
		return nil, err
	}
	if cur == nil {
		return nil, ErrRouteNotFound
	}
	pvID := cur.ProfileVersionID
	if req.ProfileVersionID != nil {
		pvID = *req.ProfileVersionID
	}
	treeID := cur.DefaultTreeNodeID
	if req.DefaultTreeNodeID != nil {
		treeID = *req.DefaultTreeNodeID
	}
	mainApp := cur.MainAppID
	if req.MainAppID != nil {
		mainApp = req.MainAppID
	}
	if err := s.validateBinding(ctx, customerID(p), pvID, treeID, mainApp); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, customerID(p), id, req, nil); err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, ErrRouteNotFound
		}
		if errors.Is(err, postgres.ErrDuplicateName) {
			return nil, ErrDuplicateRoute
		}
		return nil, err
	}
	return s.GetByID(ctx, p, id)
}

func (s *Service) QRMeta(ctx context.Context, p *platformauth.Principal, id int) (*domain.QRMeta, error) {
	if err := s.requirePerm(p); err != nil {
		return nil, err
	}
	detail, err := s.repo.GetByID(ctx, customerID(p), id)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, ErrRouteNotFound
	}
	mode := strings.TrimSpace(detail.DefaultDeviceIDMode)
	if mode == "" {
		mode = "imei"
	}
	return &domain.QRMeta{
		QRCodeKey:           detail.QRCodeKey,
		DefaultDeviceIDMode: mode,
		MainAppID:           detail.MainAppID,
	}, nil
}

func (s *Service) validateBinding(ctx context.Context, customerID, profileVersionID, treeNodeID int, mainAppID *int) error {
	if profileVersionID > 0 {
		ok, err := s.repo.IsPublishedProfileVersion(ctx, customerID, profileVersionID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrPublishedVersionRequired
		}
	}
	if treeNodeID <= 0 {
		return ErrTreeNodeRequired
	}
	ok, err := s.repo.TreeNodeBelongsToCustomer(ctx, customerID, treeNodeID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrTreeNodeRequired
	}
	if mainAppID == nil || *mainAppID <= 0 {
		return ErrMainAppRequired
	}
	return nil
}
