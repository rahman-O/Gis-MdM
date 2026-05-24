package application_test

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/application"
	"github.com/gis-mdm/server-backend-go/internal/modules/enrollment_routes/domain"
)

func TestResolveBootstrapIntent_requiresApplication(t *testing.T) {
	_, err := application.ResolveBootstrapIntent(context.Background(), nil, 1, 0, domain.BootstrapIntentStable, nil)
	if err != application.ErrMainAppRequired {
		t.Fatalf("got %v", err)
	}
}
