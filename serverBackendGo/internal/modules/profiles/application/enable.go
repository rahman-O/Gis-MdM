package application

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"

	profilepostgres "github.com/gis-mdm/server-backend-go/internal/modules/profiles/adapter/persistence/postgres"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/port"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

// EnableService toggles profile.enabled flag.
type EnableService struct {
	store port.RolloutStore
	db    *sql.DB
}

func NewEnableService(db *sql.DB) *EnableService {
	return &EnableService{store: profilepostgres.NewAssignmentRepository(db), db: db}
}

func (s *EnableService) SetEnabled(ctx context.Context, p *platformauth.Principal, profileID int, enabled bool) (*domain.EnableProfileResult, error) {
	if err := requireConfigPerm(p); err != nil {
		return nil, err
	}
	n, err := s.store.SetProfileEnabled(ctx, customerID(p), profileID, enabled)
	if err != nil {
		return nil, err
	}
	eventType := "ProfileDisabled"
	if enabled {
		eventType = "ProfileEnabled"
	}
	payload, _ := json.Marshal(map[string]any{"userId": int(p.ID)})
	_ = insertDomainEvent(ctx, s.db, eventType, strconv.Itoa(profileID), payload)
	return &domain.EnableProfileResult{Enabled: enabled, DevicesMarkedPending: n}, nil
}
