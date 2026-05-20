package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/customers/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/customers/port"
)

// CustomerRepository implements port.CustomerRepository.
type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

var _ port.CustomerRepository = (*CustomerRepository)(nil)

func (r *CustomerRepository) searchWhere(req domain.SearchRequest) (string, []any) {
	var clauses []string
	var args []any
	clauses = append(clauses, "customers.master = FALSE")
	if req.SearchValue != nil && strings.TrimSpace(*req.SearchValue) != "" {
		args = append(args, "%"+strings.TrimSpace(*req.SearchValue)+"%")
		n := len(args)
		clauses = append(clauses, fmt.Sprintf("(customers.name ILIKE $%d OR customers.description ILIKE $%d)", n, n))
	}
	if req.AccountType != nil {
		args = append(args, *req.AccountType)
		clauses = append(clauses, fmt.Sprintf("customers.accounttype = $%d", len(args)))
	}
	if req.CustomerStatus != nil && strings.TrimSpace(*req.CustomerStatus) != "" {
		st := strings.TrimSpace(*req.CustomerStatus)
		args = append(args, st)
		n := len(args)
		if st == "customer.new" {
			clauses = append(clauses, fmt.Sprintf("(customers.customerstatus = $%d OR customers.customerstatus IS NULL)", n))
		} else {
			clauses = append(clauses, fmt.Sprintf("customers.customerstatus = $%d", n))
		}
	}
	return strings.Join(clauses, " AND "), args
}

func (r *CustomerRepository) orderBy(req domain.SearchRequest) string {
	if req.SortValue != nil {
		switch *req.SortValue {
		case "registrationTime":
			if req.SortDirection != nil && strings.EqualFold(*req.SortDirection, "desc") {
				return "ORDER BY customers.registrationtime DESC NULLS LAST"
			}
			return "ORDER BY customers.registrationtime"
		case "lastLoginTime":
			if req.SortDirection != nil && strings.EqualFold(*req.SortDirection, "desc") {
				return "ORDER BY customers.lastlogintime DESC NULLS LAST"
			}
			return "ORDER BY customers.lastlogintime"
		case "expiryTime":
			if req.SortDirection != nil && strings.EqualFold(*req.SortDirection, "desc") {
				return "ORDER BY customers.expirytime DESC NULLS LAST"
			}
			return "ORDER BY customers.expirytime"
		}
	}
	return "ORDER BY customers.name"
}

func scanCustomer(row interface{ Scan(...any) error }) (*domain.Customer, error) {
	var c domain.Customer
	var id int
	var email, desc, filesDir, prefix, status sql.NullString
	var master bool
	var lastLogin, regTime, expiry sql.NullInt64
	var accountType, deviceLimit, deviceCfg sql.NullInt64
	if err := row.Scan(
		&id, &c.Name, &email, &desc, &master, &filesDir, &prefix,
		&lastLogin, &regTime, &accountType, &status, &expiry, &deviceLimit, &deviceCfg,
	); err != nil {
		return nil, err
	}
	c.ID = &id
	if email.Valid {
		c.Email = email.String
	}
	if desc.Valid {
		c.Description = desc.String
	}
	c.Master = master
	if filesDir.Valid {
		c.FilesDir = filesDir.String
	}
	if prefix.Valid {
		c.Prefix = prefix.String
	}
	if lastLogin.Valid {
		c.LastLoginTime = &lastLogin.Int64
	}
	if regTime.Valid {
		c.RegistrationTime = &regTime.Int64
	}
	if accountType.Valid {
		v := int(accountType.Int64)
		c.AccountType = &v
	}
	if status.Valid {
		c.CustomerStatus = status.String
	}
	if expiry.Valid {
		c.ExpiryTime = &expiry.Int64
	}
	if deviceLimit.Valid {
		v := int(deviceLimit.Int64)
		c.DeviceLimit = &v
	}
	if deviceCfg.Valid {
		v := int(deviceCfg.Int64)
		c.DeviceConfigurationID = &v
	}
	return &c, nil
}

const customerSelectCols = `customers.id, customers.name, customers.email, customers.description,
customers.master, customers.filesdir, customers.prefix, customers.lastlogintime,
customers.registrationtime, customers.accounttype, customers.customerstatus,
customers.expirytime, customers.devicelimit, customers.deviceconfigurationid`

func (r *CustomerRepository) Search(ctx context.Context, req domain.SearchRequest) ([]domain.Customer, error) {
	where, args := r.searchWhere(req)
	req.NormalizeSearch()
	offset := (req.CurrentPage - 1) * req.PageSize
	args = append(args, req.PageSize, offset)
	limitArg := len(args) - 1
	offsetArg := len(args)
	q := fmt.Sprintf(`SELECT %s FROM customers WHERE %s %s LIMIT $%d OFFSET $%d`,
		customerSelectCols, where, r.orderBy(req), limitArg, offsetArg)
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Customer
	for rows.Next() {
		c, err := scanCustomer(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, rows.Err()
}

func (r *CustomerRepository) Count(ctx context.Context, req domain.SearchRequest) (int64, error) {
	where, args := r.searchWhere(req)
	q := fmt.Sprintf(`SELECT COUNT(*) FROM customers WHERE %s`, where)
	var n int64
	if err := r.db.QueryRowContext(ctx, q, args...).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

func (r *CustomerRepository) GetByID(ctx context.Context, id int) (*domain.Customer, error) {
	q := fmt.Sprintf(`SELECT %s FROM customers WHERE id = $1`, customerSelectCols)
	row := r.db.QueryRowContext(ctx, q, id)
	c, err := scanCustomer(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (r *CustomerRepository) GetByName(ctx context.Context, name string) (*domain.Customer, error) {
	q := fmt.Sprintf(`SELECT %s FROM customers WHERE lower(name) = lower($1) LIMIT 1`, customerSelectCols)
	row := r.db.QueryRowContext(ctx, q, name)
	c, err := scanCustomer(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (r *CustomerRepository) GetByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	q := fmt.Sprintf(`SELECT %s FROM customers WHERE lower(email) = lower($1) LIMIT 1`, customerSelectCols)
	row := r.db.QueryRowContext(ctx, q, email)
	c, err := scanCustomer(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (r *CustomerRepository) Insert(ctx context.Context, c *domain.Customer) (int, error) {
	now := time.Now().UnixMilli()
	prefix := c.Prefix
	if prefix == "" {
		prefix = "hmdm-"
	}
	accountType := 0
	if c.AccountType != nil {
		accountType = *c.AccountType
	}
	deviceLimit := 3
	if c.DeviceLimit != nil {
		deviceLimit = *c.DeviceLimit
	}
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO customers (name, email, description, filesdir, master, prefix, registrationtime,
			accounttype, devicelimit, customerstatus, deviceconfigurationid)
		VALUES ($1,$2,$3,$4,FALSE,$5,$6,$7,$8,$9,$10) RETURNING id`,
		c.Name, nullStr(c.Email), nullStr(c.Description), c.FilesDir, prefix, now,
		accountType, deviceLimit, nullStr(c.CustomerStatus), nullInt(c.DeviceConfigurationID),
	).Scan(&id)
	return id, err
}

func (r *CustomerRepository) Update(ctx context.Context, c *domain.Customer) error {
	if c.ID == nil {
		return fmt.Errorf("customer id required")
	}
	accountType := 0
	if c.AccountType != nil {
		accountType = *c.AccountType
	}
	deviceLimit := 3
	if c.DeviceLimit != nil {
		deviceLimit = *c.DeviceLimit
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE customers SET name=$1, email=$2, description=$3, accounttype=$4, expirytime=$5,
			devicelimit=$6, customerstatus=$7, deviceconfigurationid=$8
		WHERE id=$9 AND master = FALSE`,
		c.Name, nullStr(c.Email), nullStr(c.Description), accountType, nullInt64(c.ExpiryTime),
		deviceLimit, nullStr(c.CustomerStatus), nullInt(c.DeviceConfigurationID), *c.ID,
	)
	return err
}

func (r *CustomerRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM customers WHERE id=$1 AND master = FALSE`, id)
	return err
}

func (r *CustomerRepository) PrefixUsed(ctx context.Context, prefix string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM customers WHERE lower(prefix) = lower($1))`, prefix).Scan(&exists)
	return exists, err
}

func nullStr(s string) sql.NullString {
	if strings.TrimSpace(s) == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullInt(p *int) sql.NullInt64 {
	if p == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*p), Valid: true}
}

func nullInt64(p *int64) sql.NullInt64 {
	if p == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *p, Valid: true}
}
