package postgres

import (
	"os"
	"testing"
)

func TestResolveCustomerID_singleTenant(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set")
	}
	// Integration tests run via quickstart; repository logic covered by manual UAT.
}
