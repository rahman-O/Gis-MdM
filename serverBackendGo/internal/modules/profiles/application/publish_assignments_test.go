package application_test

import (
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
)

func TestBuildImpactSummary_AssignmentBumpPreview(t *testing.T) {
	assignments := []domain.PublishImpactAssignment{
		{AssignmentID: 1, TreeNodeID: 2, TreeNodeName: "HQ", CurrentVersionNumber: 1, DeviceCount: 10},
	}
	s := application.BuildImpactSummary(5, 0, assignments)
	if len(s.AssignmentsToUpdate) != 1 {
		t.Fatalf("expected assignment row, got %+v", s.AssignmentsToUpdate)
	}
	if !s.RequiresConfirmDialog {
		t.Fatal("expected confirm when assignments present")
	}
}
