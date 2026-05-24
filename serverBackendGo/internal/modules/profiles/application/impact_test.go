package application_test

import (
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
)

func TestBuildImpactSummary_RequiresConfirmAt50(t *testing.T) {
	s := application.BuildImpactSummary(50, 1, nil)
	if !s.RequiresConfirmDialog {
		t.Fatal("expected confirm at 50 devices")
	}
	s49 := application.BuildImpactSummary(49, 1, nil)
	if s49.RequiresConfirmDialog {
		t.Fatal("expected no confirm below 50")
	}
}

func TestBuildImpactSummary_RequiresConfirmWithAssignments(t *testing.T) {
	s := application.BuildImpactSummary(0, 0, []domain.PublishImpactAssignment{{AssignmentID: 1}})
	if !s.RequiresConfirmDialog {
		t.Fatal("expected confirm when assignments will bump")
	}
}

func TestImpactConfirmThreshold(t *testing.T) {
	if application.ImpactConfirmThreshold() != 50 {
		t.Fatalf("expected threshold 50")
	}
}
