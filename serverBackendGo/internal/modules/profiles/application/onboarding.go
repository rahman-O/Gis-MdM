package application

import (
	"context"
	"database/sql"

	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// OnboardingService computes setup checklist flags for a tenant.
type OnboardingService struct {
	db *sql.DB
}

func NewOnboardingService(db *sql.DB) *OnboardingService {
	return &OnboardingService{db: db}
}

func (s *OnboardingService) Status(ctx context.Context, p *platformauth.Principal) (*domain.OnboardingStatus, error) {
	if p == nil || !p.CanManageConfigurations() {
		return nil, ErrPermissionDenied
	}
	cid := customerID(p)

	var treeBeyond bool
	_ = s.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM device_tree_nodes n
			WHERE n.customerid = $1 AND n.parent_id IS NOT NULL
		)`, cid).Scan(&treeBeyond)

	var published bool
	_ = s.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM profile_versions pv
			JOIN profiles p ON p.id = pv.profile_id
			WHERE p.customerid = $1 AND pv.status = 'published'
		)`, cid).Scan(&published)

	var routes bool
	_ = s.db.QueryRowContext(ctx, `
		SELECT EXISTS (SELECT 1 FROM enrollment_routes WHERE customerid = $1)`, cid).Scan(&routes)

	steps := []domain.OnboardingStep{
		{ID: "tree", Label: "Create a device folder in the tree", Done: treeBeyond, Path: "/devices"},
		{ID: "profile", Label: "Create and publish a profile", Done: published, Path: "/profiles"},
		{ID: "route", Label: "Create an enrollment route", Done: routes, Path: "/enrollment-routes"},
		{ID: "qr", Label: "Test enrollment QR", Done: routes, Path: "/enrollment-routes"},
	}
	complete := treeBeyond && published && routes
	return &domain.OnboardingStatus{
		Complete:            complete,
		HasTreeBeyondRoot:   treeBeyond,
		HasPublishedProfile: published,
		HasEnrollmentRoute:  routes,
		Steps:               steps,
	}, nil
}
