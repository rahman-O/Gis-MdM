package domain

// OnboardingStatus drives dashboard checklist and wizard (US6).
type OnboardingStatus struct {
	Complete             bool            `json:"complete"`
	HasTreeBeyondRoot    bool            `json:"hasTreeBeyondRoot"`
	HasPublishedProfile  bool            `json:"hasPublishedProfile"`
	HasEnrollmentRoute   bool            `json:"hasEnrollmentRoute"`
	Steps                []OnboardingStep `json:"steps"`
}

// OnboardingStep is one wizard/checklist item.
type OnboardingStep struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Done    bool   `json:"done"`
	Path    string `json:"path,omitempty"`
}
