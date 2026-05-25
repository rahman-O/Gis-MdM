package application

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	cfgdomain "github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type Service struct {
	repo           port.RouteRepository
	db             *sql.DB
	heavyThreshold int
}

func NewService(repo port.RouteRepository, db *sql.DB, heavyThreshold int) *Service {
	if heavyThreshold <= 0 {
		heavyThreshold = 500
	}
	return &Service{repo: repo, db: db, heavyThreshold: heavyThreshold}
}

var (
	ErrPermissionDenied       = errors.New("error.permission.denied")
	ErrRouteNotFound          = errors.New("error.notfound.enrollment_route")
	ErrDuplicateRoute         = errors.New("error.duplicate.enrollment_route")
	ErrTreeNodeRequired       = errors.New("error.enrollment_route.tree_node_required")
	ErrMainAppRequired        = errors.New("error.enrollment_route.main_app_required")
	ErrContainerAckRequired   = errors.New("error.enrollment_route.container_ack_required")
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

func (s *Service) List(ctx context.Context, p *platformauth.Principal) ([]domain.EnrollmentRouteView, error) {
	if err := s.requirePerm(p); err != nil {
		return nil, err
	}
	items, err := s.repo.ListViews(ctx, customerID(p))
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.EnrollmentRouteView{}
	}
	return items, nil
}

func (s *Service) GetByID(ctx context.Context, p *platformauth.Principal, id int) (*domain.EnrollmentRouteView, error) {
	if err := s.requirePerm(p); err != nil {
		return nil, err
	}
	view, err := s.repo.GetViewByID(ctx, customerID(p), id)
	if err != nil {
		return nil, err
	}
	if view == nil {
		return nil, ErrRouteNotFound
	}
	return view, nil
}

func (s *Service) ListPublishedProfileVersions(ctx context.Context, p *platformauth.Principal) ([]interface{}, error) {
	if err := s.requirePerm(p); err != nil {
		return nil, err
	}
	return []interface{}{}, nil
}

func (s *Service) Create(ctx context.Context, p *platformauth.Principal, req domain.CreateRequest) (*domain.EnrollmentRouteView, error) {
	if err := s.requirePerm(p); err != nil {
		return nil, err
	}
	cid := customerID(p)
	normalizeCreateRequest(&req)
	if err := ValidateProvisioningFields(req.WifiSSID, req.WifiPassword, req.WifiSecurityType, req.QRParameters, req.AdminExtras); err != nil {
		return nil, err
	}
	if err := s.validateCreate(ctx, cid, req); err != nil {
		return nil, err
	}
	resolved, err := ResolveBootstrapIntent(ctx, s.db, cid, req.BootstrapApplicationID, req.BootstrapIntent, req.BootstrapVersionID)
	if err != nil {
		return nil, err
	}
	key := cfgdomain.NewQRCodeKey()
	id, err := s.repo.Create(ctx, cid, req, key, resolved, req.AcknowledgeContainerPlacement)
	if err != nil {
		if errors.Is(err, postgres.ErrDuplicateName) {
			return nil, ErrDuplicateRoute
		}
		return nil, err
	}
	return s.GetByID(ctx, p, id)
}

func (s *Service) Update(ctx context.Context, p *platformauth.Principal, id int, req domain.UpdateRequest) (*domain.EnrollmentRouteView, error) {
	if err := s.requirePerm(p); err != nil {
		return nil, err
	}
	cid := customerID(p)
	cur, err := s.repo.GetViewByID(ctx, cid, id)
	if err != nil {
		return nil, err
	}
	if cur == nil {
		return nil, ErrRouteNotFound
	}
	merged := mergeUpdate(cur, req)
	if err := ValidateProvisioningFields(req.WifiSSID, req.WifiPassword, req.WifiSecurityType, req.QRParameters, req.AdminExtras); err != nil {
		return nil, err
	}
	if err := s.validateUpdate(ctx, cid, merged); err != nil {
		return nil, err
	}
	resolved, err := ResolveBootstrapIntent(ctx, s.db, cid, merged.BootstrapApplicationID, merged.BootstrapIntent, merged.BootstrapVersionID)
	if err != nil {
		return nil, err
	}
	ack := merged.AcknowledgeContainerPlacement
	if err := s.repo.Update(ctx, cid, id, toUpdateRequest(merged), &resolved, &ack); err != nil {
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
	view, err := s.repo.GetViewByID(ctx, customerID(p), id)
	if err != nil {
		return nil, err
	}
	if view == nil {
		return nil, ErrRouteNotFound
	}
	mode := strings.TrimSpace(view.DeviceIdentityMode)
	if mode == "" {
		mode = "imei"
	}
	meta := &domain.QRMeta{
		QRCodeKey:                view.QRCodeKey,
		DefaultDeviceIDMode:      mode,
		ResolvedMainAppVersionID: view.ResolvedMainAppVersionID,
		MainAppPackage:           view.ResolvedPackage,
		MainAppVersion:           view.ResolvedVersionLabel,
		TargetNodeID:             view.TargetNodeID,
		Contract: map[string]interface{}{
			"routeId":            view.ID,
			"targetNodeId":       view.TargetNodeID,
			"mainAppPackage":     view.ResolvedPackage,
			"mainAppVersion":     view.ResolvedVersionLabel,
			"deviceIdentityMode": mode,
			"bootstrapFlags":     map[string]bool{"create": true},
		},
	}
	return meta, nil
}

func normalizeCreateRequest(req *domain.CreateRequest) {
	if req.TargetNodeID <= 0 && req.DefaultTreeNodeID > 0 {
		req.TargetNodeID = req.DefaultTreeNodeID
	}
	if req.DeviceIdentityMode == nil && req.DefaultDeviceIDMode != nil {
		req.DeviceIdentityMode = req.DefaultDeviceIDMode
	}
	if strings.TrimSpace(req.BootstrapIntent) == "" {
		req.BootstrapIntent = domain.BootstrapIntentStable
	}
}

type mergedRoute struct {
	Name                          string
	Description                   string
	TargetNodeID                  int
	DeviceIdentityMode            string
	BootstrapIntent               string
	BootstrapApplicationID        int
	BootstrapVersionID            *int
	AcknowledgeContainerPlacement bool
}

func mergeUpdate(cur *domain.EnrollmentRouteView, req domain.UpdateRequest) mergedRoute {
	m := mergedRoute{
		Name:                          cur.Name,
		Description:                   cur.Description,
		TargetNodeID:                  cur.TargetNodeID,
		DeviceIdentityMode:            cur.DeviceIdentityMode,
		BootstrapIntent:               cur.BootstrapIntent,
		BootstrapApplicationID:        cur.BootstrapApplicationID,
		BootstrapVersionID:            cur.BootstrapVersionID,
		AcknowledgeContainerPlacement: cur.ContainerPlacementAcknowledged,
	}
	if req.Name != nil {
		m.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		m.Description = *req.Description
	}
	if req.TargetNodeID != nil {
		m.TargetNodeID = *req.TargetNodeID
	} else if req.DefaultTreeNodeID != nil {
		m.TargetNodeID = *req.DefaultTreeNodeID
	}
	if req.DeviceIdentityMode != nil {
		m.DeviceIdentityMode = strings.TrimSpace(*req.DeviceIdentityMode)
	} else if req.DefaultDeviceIDMode != nil {
		m.DeviceIdentityMode = strings.TrimSpace(*req.DefaultDeviceIDMode)
	}
	if req.BootstrapIntent != nil {
		m.BootstrapIntent = strings.TrimSpace(*req.BootstrapIntent)
	}
	if req.BootstrapApplicationID != nil {
		m.BootstrapApplicationID = *req.BootstrapApplicationID
	}
	if req.BootstrapVersionID != nil {
		m.BootstrapVersionID = req.BootstrapVersionID
	}
	if req.AcknowledgeContainerPlacement != nil {
		m.AcknowledgeContainerPlacement = *req.AcknowledgeContainerPlacement
	}
	return m
}

func toUpdateRequest(m mergedRoute) domain.UpdateRequest {
	name := m.Name
	desc := m.Description
	tid := m.TargetNodeID
	mode := m.DeviceIdentityMode
	intent := m.BootstrapIntent
	appID := m.BootstrapApplicationID
	ack := m.AcknowledgeContainerPlacement
	return domain.UpdateRequest{
		Name:                          &name,
		Description:                   &desc,
		TargetNodeID:                  &tid,
		DeviceIdentityMode:            &mode,
		BootstrapIntent:               &intent,
		BootstrapApplicationID:        &appID,
		BootstrapVersionID:            m.BootstrapVersionID,
		AcknowledgeContainerPlacement: &ack,
	}
}

func (s *Service) validateCreate(ctx context.Context, customerID int, req domain.CreateRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrTreeNodeRequired
	}
	if req.TargetNodeID <= 0 {
		return ErrTreeNodeRequired
	}
	if req.BootstrapApplicationID <= 0 {
		return ErrMainAppRequired
	}
	return s.validateNode(ctx, customerID, req.TargetNodeID, req.AcknowledgeContainerPlacement)
}

func (s *Service) validateUpdate(ctx context.Context, customerID int, m mergedRoute) error {
	if strings.TrimSpace(m.Name) == "" {
		return ErrTreeNodeRequired
	}
	if m.TargetNodeID <= 0 {
		return ErrTreeNodeRequired
	}
	if m.BootstrapApplicationID <= 0 {
		return ErrMainAppRequired
	}
	return s.validateNode(ctx, customerID, m.TargetNodeID, m.AcknowledgeContainerPlacement)
}

func (s *Service) validateNode(ctx context.Context, customerID, nodeID int, ack bool) error {
	ok, err := s.repo.TreeNodeBelongsToCustomer(ctx, customerID, nodeID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrTreeNodeRequired
	}
	kind, err := s.repo.NodePlacementKind(ctx, customerID, nodeID)
	if err != nil {
		return err
	}
	if kind == "inheritable" && !ack {
		return ErrContainerAckRequired
	}
	return nil
}
