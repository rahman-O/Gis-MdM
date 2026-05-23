package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/audit/domain"
)

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Search(ctx context.Context, customerID int64, f domain.AuditLogFilter) ([]domain.AuditLogRecord, int64, error) {
	if f.PageSize <= 0 {
		f.PageSize = 50
	}
	if f.PageNum <= 0 {
		f.PageNum = 1
	}
	where := []string{"customerid = $1"}
	args := []any{customerID}
	n := 2
	if f.UserFilter != "" {
		where = append(where, fmt.Sprintf("login ILIKE $%d", n))
		args = append(args, "%"+f.UserFilter+"%")
		n++
	}
	if f.MessageFilter != "" {
		where = append(where, fmt.Sprintf("(ipaddress ILIKE $%d OR login ILIKE $%d OR action ILIKE $%d)", n, n, n))
		args = append(args, "%"+f.MessageFilter+"%")
		n++
	}
	if f.DateFrom != nil {
		where = append(where, fmt.Sprintf("createtime >= $%d", n))
		args = append(args, *f.DateFrom)
		n++
	}
	if f.DateTo != nil {
		where = append(where, fmt.Sprintf("createtime <= $%d", n))
		args = append(args, *f.DateTo)
		n++
	}
	w := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM plugin_audit_log WHERE `+w, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (f.PageNum - 1) * f.PageSize
	q := `SELECT id, createtime, customerid, userid, COALESCE(login,''), COALESCE(action,''),
		COALESCE(payload,''), COALESCE(ipaddress,''), errorcode
		FROM plugin_audit_log WHERE ` + w + ` ORDER BY createtime DESC LIMIT $` + fmt.Sprint(n) + ` OFFSET $` + fmt.Sprint(n+1)
	args = append(args, f.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []domain.AuditLogRecord
	for rows.Next() {
		var rec domain.AuditLogRecord
		var uid sql.NullInt64
		if err := rows.Scan(&rec.ID, &rec.CreateTime, &rec.CustomerID, &uid, &rec.Login, &rec.Action, &rec.Payload, &rec.IPAddress, &rec.ErrorCode); err != nil {
			return nil, 0, err
		}
		if uid.Valid {
			v := uid.Int64
			rec.UserID = &v
		}
		items = append(items, rec)
	}
	return items, total, rows.Err()
}
