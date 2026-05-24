package domain

// EnrollmentDeleteImpact is returned by GET /enrollment-routes/:id/impact.
type EnrollmentDeleteImpact struct {
	EnrollingNowCount        int `json:"enrollingNowCount"`
	HistoricalEnrolledCount  int `json:"historicalEnrolledCount"`
	ActiveQrScans7d          int `json:"activeQrScans7d"`
}
