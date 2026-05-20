package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Connect opens a PostgreSQL connection. Returns nil, nil when DATABASE_URL is empty (scaffold mode).
func Connect(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, nil
	}
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}
