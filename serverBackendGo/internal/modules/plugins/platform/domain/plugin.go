package domain

// Plugin catalog row (JSON matches React PluginRow).
type Plugin struct {
	ID                   int64  `json:"id"`
	Identifier           string `json:"identifier,omitempty"`
	Name                 string `json:"name,omitempty"`
	NameLocalizationKey  string `json:"nameLocalizationKey,omitempty"`
	Description          string `json:"description,omitempty"`
	Disabled             bool   `json:"disabled,omitempty"`
}
