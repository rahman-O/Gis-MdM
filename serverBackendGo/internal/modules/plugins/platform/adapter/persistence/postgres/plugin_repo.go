package postgres

import (
	"context"
	"database/sql"
	"strings"

	"github.com/lib/pq"

	"github.com/gis-mdm/server-backend-go/internal/config"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/platform/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/platform/port"
)

type PluginRepository struct {
	db     *sql.DB
	enabled func(string) bool
}

func NewPluginRepository(db *sql.DB, cfg config.Config) *PluginRepository {
	return &PluginRepository{
		db: db,
		enabled: cfg.IsPluginEnabled,
	}
}

var _ port.Repository = (*PluginRepository)(nil)

func (r *PluginRepository) filterEnabled(rows []domain.Plugin) []domain.Plugin {
	out := make([]domain.Plugin, 0, len(rows))
	for _, p := range rows {
		if r.enabled(p.Identifier) {
			out = append(out, p)
		}
	}
	return out
}

func (r *PluginRepository) scanPlugins(rows *sql.Rows) ([]domain.Plugin, error) {
	defer rows.Close()
	var list []domain.Plugin
	for rows.Next() {
		var p domain.Plugin
		var desc, js, fn, st, sk sql.NullString
		if err := rows.Scan(&p.ID, &p.Identifier, &p.Name, &desc, &p.Disabled, &js, &fn, &st, &sk); err != nil {
			return nil, err
		}
		if sk.Valid {
			p.NameLocalizationKey = sk.String
		}
		if desc.Valid {
			p.Description = desc.String
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

const pluginSelect = `SELECT id, identifier, name, description, disabled,
	javascriptmodulefile, functionsviewtemplate, settingsviewtemplate, namelocalizationkey
	FROM plugins`

func (r *PluginRepository) FindActive(ctx context.Context) ([]domain.Plugin, error) {
	rows, err := r.db.QueryContext(ctx, pluginSelect+` WHERE disabled = FALSE ORDER BY identifier`)
	if err != nil {
		return nil, err
	}
	list, err := r.scanPlugins(rows)
	if err != nil {
		return nil, err
	}
	return r.filterEnabled(list), nil
}

func (r *PluginRepository) FindAvailableByCustomer(ctx context.Context, customerID int64) ([]domain.Plugin, error) {
	rows, err := r.db.QueryContext(ctx, pluginSelect+`
		WHERE disabled = FALSE
		AND NOT EXISTS (
			SELECT 1 FROM pluginsdisabled pd
			WHERE pd.pluginid = plugins.id AND pd.customerid = $1
		)
		ORDER BY identifier`, customerID)
	if err != nil {
		return nil, err
	}
	list, err := r.scanPlugins(rows)
	if err != nil {
		return nil, err
	}
	return r.filterEnabled(list), nil
}

func (r *PluginRepository) FindRegistered(ctx context.Context) ([]domain.Plugin, error) {
	rows, err := r.db.QueryContext(ctx, pluginSelect+` ORDER BY identifier`)
	if err != nil {
		return nil, err
	}
	list, err := r.scanPlugins(rows)
	if err != nil {
		return nil, err
	}
	return r.filterEnabled(list), nil
}

func (r *PluginRepository) SaveDisabled(ctx context.Context, customerID int64, pluginIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM pluginsdisabled WHERE customerid = $1`, customerID); err != nil {
		return err
	}
	if len(pluginIDs) > 0 {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO pluginsdisabled (pluginid, customerid)
			SELECT unnest($1::int[]), $2`, pq.Array(pluginIDs), customerID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// PluginIDByIdentifier resolves catalog id (for status cache).
func PluginIDByIdentifier(ctx context.Context, db *sql.DB, identifier string) (int64, error) {
	var id int64
	err := db.QueryRowContext(ctx, `SELECT id FROM plugins WHERE lower(identifier) = lower($1)`, identifier).Scan(&id)
	return id, err
}

// IdentifierForPluginIDs returns map id->identifier.
func IdentifierForPluginIDs(ctx context.Context, db *sql.DB, ids []int64) (map[int64]string, error) {
	if len(ids) == 0 {
		return map[int64]string{}, nil
	}
	rows, err := db.QueryContext(ctx, `SELECT id, identifier FROM plugins WHERE id = ANY($1)`, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := make(map[int64]string)
	for rows.Next() {
		var id int64
		var ident string
		if err := rows.Scan(&id, &ident); err != nil {
			return nil, err
		}
		m[id] = strings.ToLower(ident)
	}
	return m, rows.Err()
}
