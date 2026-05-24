package application_test

import (
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
)

func TestComputeHealth_draftOnly(t *testing.T) {
	h, _, badges := application.ComputeHealth(domain.HubMetrics{HasPublished: false})
	if h != domain.HealthDraftOnly {
		t.Fatalf("got %s", h)
	}
	if len(badges) != 1 || badges[0] != "draft_only" {
		t.Fatalf("badges %v", badges)
	}
}

func TestComputeHealth_healthy(t *testing.T) {
	h, reasons, _ := application.ComputeHealth(domain.HubMetrics{
		HasPublished: true, Enabled: true, AssignmentCount: 2,
	})
	if h != domain.HealthHealthy {
		t.Fatalf("got %s reasons %v", h, reasons)
	}
}

func TestComputeHealth_noAssignment(t *testing.T) {
	h, _, badges := application.ComputeHealth(domain.HubMetrics{
		HasPublished: true, Enabled: true, AssignmentCount: 0,
	})
	if h != domain.HealthWarning {
		t.Fatalf("got %s", h)
	}
	found := false
	for _, b := range badges {
		if b == "no_assignment" {
			found = true
		}
	}
	if !found {
		t.Fatalf("badges %v", badges)
	}
}

func TestComputeHealth_rolloutError(t *testing.T) {
	h, _, _ := application.ComputeHealth(domain.HubMetrics{
		HasPublished: true, Enabled: true, AssignmentCount: 1, RolloutFailureCount: 3,
	})
	if h != domain.HealthError {
		t.Fatalf("got %s", h)
	}
}
