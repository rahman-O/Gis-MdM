package application

import "github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"

const impactConfirmThreshold = 50

// BuildImpactSummary maps device/route counts and assignment rows to API impact response.
func BuildImpactSummary(deviceCount, routeCount int, assignments []domain.PublishImpactAssignment) domain.ImpactSummary {
	if assignments == nil {
		assignments = []domain.PublishImpactAssignment{}
	}
	requires := deviceCount >= impactConfirmThreshold || len(assignments) > 0
	return domain.ImpactSummary{
		DeviceCount:             deviceCount,
		EnrollmentRouteCount:    routeCount,
		RequiresConfirmDialog:   requires,
		AssignmentsToUpdate:     assignments,
	}
}

// ImpactConfirmThreshold exposes the FR-006 threshold for tests.
func ImpactConfirmThreshold() int { return impactConfirmThreshold }
