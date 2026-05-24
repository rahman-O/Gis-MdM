package application

import "testing"

func TestEnrollmentStateConstants(t *testing.T) {
	if EnrollmentStateEnrolled == "" || EnrollmentStateActive == "" {
		t.Fatal("enrollment state constants must be non-empty")
	}
}
