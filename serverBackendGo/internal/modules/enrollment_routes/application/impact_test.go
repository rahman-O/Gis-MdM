package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// impactRepo is a test double for the RouteRepository focused on Impact tests.
type impactRepo struct {
	fakeRepo
	view   *domain.EnrollmentRouteView
	impact *domain.EnrollmentDeleteImpact
	err    error
}

func (r *impactRepo) GetViewByID(_ context.Context, _ int, _ int) (*domain.EnrollmentRouteView, error) {
	return r.view, r.err
}

func (r *impactRepo) DeleteImpact(_ context.Context, _ int, _ int) (*domain.EnrollmentDeleteImpact, error) {
	if r.impact == nil && r.err == nil {
		return &domain.EnrollmentDeleteImpact{}, nil
	}
	return r.impact, r.err
}

func impactPrincipal() *platformauth.Principal {
	return &platformauth.Principal{CustomerID: 1, Permissions: []string{"configurations"}}
}

func TestImpact(t *testing.T) {
	tests := []struct {
		name      string
		principal *platformauth.Principal
		repo      *impactRepo
		wantErr   error
		want      *domain.EnrollmentDeleteImpact
	}{
		{
			name:      "permission denied when principal is nil",
			principal: nil,
			repo: &impactRepo{
				view: &domain.EnrollmentRouteView{ID: 1},
			},
			wantErr: application.ErrPermissionDenied,
		},
		{
			name:      "permission denied when principal lacks permission",
			principal: &platformauth.Principal{CustomerID: 1, Permissions: []string{}},
			repo: &impactRepo{
				view: &domain.EnrollmentRouteView{ID: 1},
			},
			wantErr: application.ErrPermissionDenied,
		},
		{
			name:      "route not found when view is nil",
			principal: impactPrincipal(),
			repo: &impactRepo{
				view: nil,
			},
			wantErr: application.ErrRouteNotFound,
		},
		{
			name:      "zero counts",
			principal: impactPrincipal(),
			repo: &impactRepo{
				view: &domain.EnrollmentRouteView{ID: 1},
				impact: &domain.EnrollmentDeleteImpact{
					EnrollingNowCount:       0,
					HistoricalEnrolledCount: 0,
					ActiveQrScans7d:         0,
				},
			},
			want: &domain.EnrollmentDeleteImpact{
				EnrollingNowCount:       0,
				HistoricalEnrolledCount: 0,
				ActiveQrScans7d:         0,
			},
		},
		{
			name:      "non-zero counts",
			principal: impactPrincipal(),
			repo: &impactRepo{
				view: &domain.EnrollmentRouteView{ID: 5},
				impact: &domain.EnrollmentDeleteImpact{
					EnrollingNowCount:       3,
					HistoricalEnrolledCount: 42,
					ActiveQrScans7d:         7,
				},
			},
			want: &domain.EnrollmentDeleteImpact{
				EnrollingNowCount:       3,
				HistoricalEnrolledCount: 42,
				ActiveQrScans7d:         7,
			},
		},
		{
			name:      "large counts",
			principal: impactPrincipal(),
			repo: &impactRepo{
				view: &domain.EnrollmentRouteView{ID: 10},
				impact: &domain.EnrollmentDeleteImpact{
					EnrollingNowCount:       1000,
					HistoricalEnrolledCount: 50000,
					ActiveQrScans7d:         9999,
				},
			},
			want: &domain.EnrollmentDeleteImpact{
				EnrollingNowCount:       1000,
				HistoricalEnrolledCount: 50000,
				ActiveQrScans7d:         9999,
			},
		},
		{
			name:      "repo error propagates",
			principal: impactPrincipal(),
			repo: &impactRepo{
				view: nil,
				err:  errors.New("db connection failed"),
			},
			wantErr: errors.New("db connection failed"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := application.NewService(tc.repo, nil, 500)
			got, err := svc.Impact(context.Background(), tc.principal, 1)

			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tc.wantErr)
				}
				if err.Error() != tc.wantErr.Error() {
					t.Fatalf("expected error %q, got %q", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("expected non-nil impact result")
			}
			if got.EnrollingNowCount != tc.want.EnrollingNowCount {
				t.Errorf("EnrollingNowCount: got %d, want %d", got.EnrollingNowCount, tc.want.EnrollingNowCount)
			}
			if got.HistoricalEnrolledCount != tc.want.HistoricalEnrolledCount {
				t.Errorf("HistoricalEnrolledCount: got %d, want %d", got.HistoricalEnrolledCount, tc.want.HistoricalEnrolledCount)
			}
			if got.ActiveQrScans7d != tc.want.ActiveQrScans7d {
				t.Errorf("ActiveQrScans7d: got %d, want %d", got.ActiveQrScans7d, tc.want.ActiveQrScans7d)
			}
		})
	}
}
