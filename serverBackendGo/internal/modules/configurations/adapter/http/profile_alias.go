package http

import (
	"context"
	"database/sql"

	cfgdomain "github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	profilepostgres "github.com/gis-mdm/server-backend-go/internal/modules/profiles/adapter/persistence/postgres"
)

// ProfileAlias resolves legacy configuration IDs to profile version payloads (transition T087).
type ProfileAlias struct {
	db   *sql.DB
	repo *profilepostgres.ProfileRepository
}

func NewProfileAlias(db *sql.DB) *ProfileAlias {
	repo := profilepostgres.NewProfileRepository(db)
	return &ProfileAlias{db: db, repo: repo}
}

func (a *ProfileAlias) GetByConfigurationID(ctx context.Context, customerID, configurationID int) (map[string]any, bool, error) {
	profileID, versionID, ok, err := a.resolveVersion(ctx, customerID, configurationID)
	if err != nil || !ok {
		return nil, false, err
	}
	payload, meta, err := a.repo.GetVersion(ctx, customerID, profileID, versionID)
	if err != nil || payload == nil || meta == nil {
		return nil, false, err
	}
	data := cfgdomain.ConfigurationResponseMap(payload)
	data["id"] = configurationID
	data["profileId"] = profileID
	data["versionId"] = meta.VersionID
	return data, true, nil
}

func (a *ProfileAlias) SaveByConfigurationID(ctx context.Context, customerID, configurationID int, payload cfgdomain.Configuration) (map[string]any, bool, error) {
	profileID, versionID, ok, err := a.resolveVersion(ctx, customerID, configurationID)
	if err != nil || !ok {
		return nil, false, err
	}
	if versionID <= 0 {
		versionID, err = a.repo.EnsureDraft(ctx, customerID, profileID)
		if err != nil {
			return nil, true, err
		}
	}
	id := profileID
	payload.ID = &id
	if err := a.repo.SaveDraft(ctx, customerID, profileID, versionID, payload); err != nil {
		return nil, true, err
	}
	return a.GetByConfigurationID(ctx, customerID, configurationID)
}

func (a *ProfileAlias) resolveVersion(ctx context.Context, customerID, configurationID int) (profileID, versionID int, ok bool, err error) {
	var pubID, draftID sql.NullInt64
	err = a.db.QueryRowContext(ctx, `
		SELECT p.id, p.published_version_id, p.draft_version_id
		FROM profiles p
		WHERE p.customerid = $1 AND p.legacy_configuration_id = $2`,
		customerID, configurationID).Scan(&profileID, &pubID, &draftID)
	if err == sql.ErrNoRows {
		return 0, 0, false, nil
	}
	if err != nil {
		return 0, 0, false, err
	}
	if draftID.Valid && draftID.Int64 > 0 {
		return profileID, int(draftID.Int64), true, nil
	}
	if pubID.Valid && pubID.Int64 > 0 {
		return profileID, int(pubID.Int64), true, nil
	}
	return profileID, 0, true, nil
}
