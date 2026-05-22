package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
)

// DeviceStatusRepository upserts devicestatuses from agent info JSON.
type DeviceStatusRepository struct {
	db *sql.DB
}

func NewDeviceStatusRepository(db *sql.DB) *DeviceStatusRepository {
	return &DeviceStatusRepository{db: db}
}

type agentInfoPayload struct {
	Applications []struct {
		Pkg     string `json:"pkg"`
		Version string `json:"version"`
	} `json:"applications"`
	Files []struct {
		Path string `json:"path"`
	} `json:"files"`
}

// UpsertFromInfoJSON derives install statuses and upserts devicestatuses (simplified Java parity).
func (r *DeviceStatusRepository) UpsertFromInfoJSON(ctx context.Context, deviceID int, infoJSON string) error {
	appsStatus := "FAILURE"
	filesStatus := "OTHER"
	trimmed := strings.TrimSpace(infoJSON)
	if trimmed != "" {
		var payload agentInfoPayload
		if err := json.Unmarshal([]byte(trimmed), &payload); err == nil {
			if len(payload.Applications) > 0 {
				appsStatus = "SUCCESS"
			}
			if len(payload.Files) > 0 {
				filesStatus = "SUCCESS"
			}
		}
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO devicestatuses (deviceid, applicationsstatus, configfilesstatus)
		VALUES ($1, $2, $3)
		ON CONFLICT (deviceid) DO UPDATE SET
			applicationsstatus = EXCLUDED.applicationsstatus,
			configfilesstatus = EXCLUDED.configfilesstatus`,
		deviceID, appsStatus, filesStatus)
	return err
}
