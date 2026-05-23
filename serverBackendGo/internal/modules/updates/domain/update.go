package domain

// UpdateEntry mirrors Java UpdateEntry for the control panel.
type UpdateEntry struct {
	Pkg                string `json:"pkg"`
	Version            string `json:"version"`
	CurrentVersion     string `json:"currentVersion,omitempty"`
	URL                string `json:"url,omitempty"`
	Description        string `json:"description,omitempty"`
	Outdated           bool   `json:"outdated,omitempty"`
	Downloaded         bool   `json:"downloaded,omitempty"`
	UpdateDisabled     bool   `json:"updateDisabled,omitempty"`
	UpdateDisableReason string `json:"updateDisableReason,omitempty"`
}

// UpdateRequest is POST /private/update body.
type UpdateRequest struct {
	Updates   []UpdateEntry `json:"updates"`
	Update    bool          `json:"update"`
	SendStats bool          `json:"sendStats"`
}
