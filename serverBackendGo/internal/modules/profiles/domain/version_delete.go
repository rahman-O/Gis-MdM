package domain

// VersionDeleteResult is returned by DELETE /profiles/:id/versions/:versionId.
type VersionDeleteResult struct {
	ProfileID int `json:"profileId"`
	VersionID int `json:"versionId"`
}
