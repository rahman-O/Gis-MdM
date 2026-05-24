package port

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
)

// RolloutStore reads/writes assignments and device rollout fields.
type RolloutStore interface {
	CountSubtreeDevices(ctx context.Context, customerID, treeNodeID int) (int, error)
	ListAssignments(ctx context.Context, customerID, profileID int) ([]domain.ProfileTreeAssignment, error)
	GetAssignmentImpact(ctx context.Context, customerID, treeNodeID int) (deviceCount int, folderName string, err error)
	UpsertAssignment(ctx context.Context, customerID, profileID, profileVersionID, treeNodeID, createdBy int) (assignmentID int, err error)
	DeleteAssignment(ctx context.Context, customerID, profileID, assignmentID int) error
	MarkSubtreePending(ctx context.Context, customerID, treeNodeID, targetVersionID int) (affected int, err error)
	ListRolloutDevices(ctx context.Context, customerID, profileID int, q domain.RolloutDevicesQuery) ([]domain.DeviceRolloutRow, int, error)
	IsVersionPublished(ctx context.Context, customerID, profileID, versionID int) (bool, error)
	IsProfileEnabled(ctx context.Context, customerID, profileID int) (bool, error)
	SetProfileEnabled(ctx context.Context, customerID, profileID int, enabled bool) (devicesPending int, err error)
	MarkProfileDevicesPending(ctx context.Context, customerID, profileID int) (int, error)
	ListAssignmentsForPublishImpact(ctx context.Context, customerID, profileID int) ([]domain.PublishImpactAssignment, error)
	BumpAllAssignmentsOnPublish(ctx context.Context, customerID, profileID, newVersionID int) (assignmentsUpdated, devicesAffected int, err error)
}
