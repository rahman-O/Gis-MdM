package postgres

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CustomerFilesRepo resolves customer filesdir for uploads.
type CustomerFilesRepo struct {
	db *sql.DB
}

func NewCustomerFilesRepo(db *sql.DB) *CustomerFilesRepo {
	return &CustomerFilesRepo{db: db}
}

// FilesDir returns filesdir for a tenant.
func (r *CustomerFilesRepo) FilesDir(c *gin.Context, customerID int) (string, error) {
	var dir sql.NullString
	err := r.db.QueryRowContext(c.Request.Context(), `
		SELECT filesdir FROM customers WHERE id = $1`, customerID).Scan(&dir)
	if err != nil {
		return "", err
	}
	if !dir.Valid || dir.String == "" {
		return "customer-" + strconv.Itoa(customerID), nil
	}
	return dir.String, nil
}
