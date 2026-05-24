package application

import (
	"context"
	"errors"

	"github.com/gis-mdm/server-backend-go/internal/modules/device_tree/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/device_tree/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// Service implements device tree use cases.
type Service struct {
	repo port.TreeRepository
}

func NewService(repo port.TreeRepository) *Service {
	return &Service{repo: repo}
}

var ErrPermissionDenied = errors.New("error.permission.denied")

func customerID(p *platformauth.Principal) int {
	if p == nil {
		return 0
	}
	return p.CustomerID
}

func (s *Service) List(ctx context.Context, p *platformauth.Principal) (*domain.TreeList, error) {
	if p == nil {
		return nil, ErrPermissionDenied
	}
	cid := customerID(p)
	rootID, err := s.repo.EnsureRoot(ctx, cid)
	if err != nil {
		return nil, err
	}
	nodes, err := s.repo.ListByCustomer(ctx, cid)
	if err != nil {
		return nil, err
	}
	return &domain.TreeList{Nodes: nodes, RootID: rootID}, nil
}

func (s *Service) Create(ctx context.Context, p *platformauth.Principal, req domain.CreateNodeRequest) (*domain.TreeNode, error) {
	if p == nil {
		return nil, ErrPermissionDenied
	}
	return s.repo.Create(ctx, customerID(p), req)
}

func (s *Service) Update(ctx context.Context, p *platformauth.Principal, id int, req domain.UpdateNodeRequest) (*domain.TreeNode, error) {
	if p == nil {
		return nil, ErrPermissionDenied
	}
	return s.repo.Update(ctx, customerID(p), id, req)
}

func (s *Service) Delete(ctx context.Context, p *platformauth.Principal, id int, req domain.DeleteNodeRequest) error {
	if p == nil {
		return ErrPermissionDenied
	}
	cid := customerID(p)
	count, err := s.repo.CountDevicesInSubtree(ctx, cid, id)
	if err != nil {
		return err
	}
	if count > 0 {
		if req.TargetNodeID <= 0 {
			return domain.ErrTargetRequired
		}
		return s.repo.DeleteWithRelocation(ctx, cid, id, req.TargetNodeID)
	}
	// No devices: still require valid target if children exist — delete subtree via relocation to parent or explicit target
	if req.TargetNodeID <= 0 {
		node, err := s.repo.GetByID(ctx, cid, id)
		if err != nil {
			return err
		}
		if node.ParentID == nil {
			return domain.ErrCannotDeleteRoot
		}
		return s.repo.DeleteWithRelocation(ctx, cid, id, *node.ParentID)
	}
	return s.repo.DeleteWithRelocation(ctx, cid, id, req.TargetNodeID)
}
