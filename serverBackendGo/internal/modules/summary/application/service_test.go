package application

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/summary/domain"
)

type stubRepo struct{}

func (stubRepo) HasDevicesTable(context.Context) (bool, error) { return false, nil }
func (stubRepo) GetDeviceStats(context.Context, int, int64) (*domain.DeviceStats, error) {
	return domain.EmptyDeviceStats(), nil
}

func TestGetDeviceStats_emptyShape(t *testing.T) {
	svc := NewService(stubRepo{})
	stats, err := svc.GetDeviceStats(context.Background(), 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(stats.StatusSummary) != 3 {
		t.Fatalf("statusSummary len=%d", len(stats.StatusSummary))
	}
	if stats.StatusSummary[0].StringAttr != "green" {
		t.Fatalf("expected green first, got %q", stats.StatusSummary[0].StringAttr)
	}
	if len(stats.DevicesEnrolledMonthly) != 12 {
		t.Fatalf("monthly len=%d", len(stats.DevicesEnrolledMonthly))
	}
}
