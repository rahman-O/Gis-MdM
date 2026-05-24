package domain

// Bootstrap intent modes (021).
const (
	BootstrapIntentStable   = "stable"
	BootstrapIntentSpecific = "specific"
	BootstrapIntentLatest   = "latest"
)

// ValidBootstrapIntent reports whether s is a known intent.
func ValidBootstrapIntent(s string) bool {
	switch s {
	case BootstrapIntentStable, BootstrapIntentSpecific, BootstrapIntentLatest:
		return true
	default:
		return false
	}
}
