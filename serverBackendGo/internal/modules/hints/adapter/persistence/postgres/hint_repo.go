package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gis-mdm/server-backend-go/internal/modules/hints/port"
)

// HintRepository implements hints persistence.
type HintRepository struct {
	db *sql.DB
}

func NewHintRepository(db *sql.DB) *HintRepository {
	return &HintRepository{db: db}
}

var _ port.Repository = (*HintRepository)(nil)

func (r *HintRepository) GetHistory(ctx context.Context, userID int64) ([]string, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT hintkey FROM userhints WHERE userid = $1 ORDER BY hintkey`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("get hint history: %w", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		out = append(out, key)
	}
	return out, rows.Err()
}

func (r *HintRepository) MarkShown(ctx context.Context, userID int64, hintKey string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO userhints (userid, hintkey) VALUES ($1, $2)
		ON CONFLICT (userid, hintkey) DO NOTHING`,
		userID, hintKey,
	)
	return err
}

func (r *HintRepository) Enable(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM userhints WHERE userid = $1`, userID)
	return err
}

func (r *HintRepository) Disable(ctx context.Context, userID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM userhints WHERE userid = $1`, userID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO userhints (userid, hintkey)
		SELECT $1, hintkey FROM userhinttypes`,
		userID,
	); err != nil {
		return err
	}
	return tx.Commit()
}
