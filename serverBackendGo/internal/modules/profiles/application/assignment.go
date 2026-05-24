package application

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	profilepostgres "github.com/gis-mdm/server-backend-go/internal/modules/profiles/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

var (
	ErrProfileDisabled          = errors.New("error.profile.disabled")
	ErrVersionNotPublished      = errors.New("error.profile.version.notPublished")
	ErrAssignmentConfirmRequired = errors.New("error.profile.assignment.confirmRequired")
	ErrAssignmentNotFound       = errors.New("error.profile.assignment.nodeNotFound")
)

// AssignmentService manages tree folder assignments.
type AssignmentService struct {
	store port.RolloutStore
	db    *sql.DB
}

func NewAssignmentService(db *sql.DB) *AssignmentService {
	return &AssignmentService{store: profilepostgres.NewAssignmentRepository(db), db: db}
}

func (s *AssignmentService) List(ctx context.Context, p *platformauth.Principal, profileID int) ([]domain.ProfileTreeAssignment, error) {
	if err := requireConfigPerm(p); err != nil {
		return nil, err
	}
	return s.store.ListAssignments(ctx, customerID(p), profileID)
}

func (s *AssignmentService) Impact(ctx context.Context, p *platformauth.Principal, profileID, treeNodeID int) (*domain.AssignmentImpact, error) {
	if err := requireConfigPerm(p); err != nil {
		return nil, err
	}
	n, name, err := s.store.GetAssignmentImpact(ctx, customerID(p), treeNodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAssignmentNotFound
		}
		return nil, err
	}
	return &domain.AssignmentImpact{
		DeviceCount:           n,
		RequiresConfirmDialog: n >= ImpactConfirmThreshold(),
		FolderName:            name,
	}, nil
}

func (s *AssignmentService) Put(ctx context.Context, p *platformauth.Principal, profileID int, req domain.PutAssignmentRequest) (*domain.PutAssignmentResult, error) {
	if err := requireConfigPerm(p); err != nil {
		return nil, err
	}
	cid := customerID(p)
	enabled, err := s.store.IsProfileEnabled(ctx, cid, profileID)
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, ErrProfileDisabled
	}
	pub, err := s.store.IsVersionPublished(ctx, cid, profileID, req.ProfileVersionID)
	if err != nil {
		return nil, err
	}
	if !pub {
		return nil, ErrVersionNotPublished
	}
	impact, err := s.Impact(ctx, p, profileID, req.TreeNodeID)
	if err != nil {
		return nil, err
	}
	if impact.RequiresConfirmDialog && !req.ConfirmImpact {
		return nil, ErrAssignmentConfirmRequired
	}
	assignmentID, err := s.store.UpsertAssignment(ctx, cid, profileID, req.ProfileVersionID, req.TreeNodeID, int(p.ID))
	if err != nil {
		return nil, err
	}
	affected, err := s.store.MarkSubtreePending(ctx, cid, req.TreeNodeID, req.ProfileVersionID)
	if err != nil {
		return nil, err
	}
	payload, _ := json.Marshal(map[string]any{
		"profileId": profileID, "treeNodeId": req.TreeNodeID, "versionId": req.ProfileVersionID,
	})
	_ = insertDomainEvent(ctx, s.db, "ProfileAssignmentChanged", fmt.Sprintf("profile:%d", profileID), payload)
	list, err := s.store.ListAssignments(ctx, cid, profileID)
	if err != nil {
		return nil, err
	}
	var item domain.ProfileTreeAssignment
	for _, a := range list {
		if a.AssignmentID == assignmentID {
			item = a
			break
		}
	}
	return &domain.PutAssignmentResult{ProfileTreeAssignment: item, AffectedDevices: affected}, nil
}

func (s *AssignmentService) Delete(ctx context.Context, p *platformauth.Principal, profileID, assignmentID int) error {
	if err := requireConfigPerm(p); err != nil {
		return err
	}
	err := s.store.DeleteAssignment(ctx, customerID(p), profileID, assignmentID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrAssignmentNotFound
	}
	return err
}

func requireConfigPerm(p *platformauth.Principal) error {
	if p == nil || !p.CanManageConfigurations() {
		return ErrPermissionDenied
	}
	return nil
}

func insertDomainEvent(ctx context.Context, db *sql.DB, eventType, aggregateID string, payload []byte) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO domain_events (event_type, aggregate_id, payload)
		VALUES ($1, $2, $3)`, eventType, aggregateID, payload)
	return err
}
