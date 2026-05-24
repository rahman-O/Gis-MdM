package application_test

import (
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
)

func TestEffectiveProfileResolutionZeroValue(t *testing.T) {
	r := domain.EffectiveProfileResolution{}
	if r.Source != "" {
		t.Fatalf("expected empty source on zero value")
	}
	if r.ProfileVersionID != 0 {
		t.Fatalf("expected no version on zero value")
	}
}
