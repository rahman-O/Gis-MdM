package application

import "github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"

// ComputeHealth derives health, reasons, and list badges from metrics.
func ComputeHealth(m domain.HubMetrics) (health string, reasons, badges []string) {
	reasons = []string{}
	badges = []string{}

	if !m.HasPublished {
		return domain.HealthDraftOnly, []string{"no_published"}, []string{"draft_only"}
	}

	if m.AssignmentCount == 0 {
		reasons = append(reasons, "no_assignment")
		badges = append(badges, "no_assignment")
	}
	if m.HasUnpublishedDraft {
		badges = append(badges, "draft_changes")
	}
	if m.StalePublish {
		reasons = append(reasons, "stale_publish")
		badges = append(badges, "stale")
	}
	if m.RolloutFailureCount > 0 {
		reasons = append(reasons, "rollout_failures")
		badges = append(badges, "rollout_issues")
	}
	if !m.Enabled {
		badges = append(badges, "disabled")
	}

	health = domain.HealthHealthy
	if m.RolloutFailureCount > 0 && m.Enabled {
		health = domain.HealthError
		reasons = append(reasons, "rollout_error")
	} else if len(reasons) > 0 || !m.Enabled {
		health = domain.HealthWarning
		if !m.Enabled && m.RolloutFailureCount > 0 {
			health = domain.HealthError
		}
	}
	return health, reasons, badges
}

// LifecycleLabel returns draft | published | disabled for UI.
func LifecycleLabel(enabled, hasPublished bool) string {
	if !enabled {
		return "disabled"
	}
	if !hasPublished {
		return "draft"
	}
	return "published"
}
