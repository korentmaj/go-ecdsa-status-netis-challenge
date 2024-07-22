package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(connStr string) error {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open the database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Create tables
	createTables := `
	CREATE TABLE IF NOT EXISTS status_lists (
		id SERIAL PRIMARY KEY,
		encoded_list BYTEA NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS statuses (
		id SERIAL PRIMARY KEY,
		list_id INTEGER REFERENCES status_lists(id) ON DELETE CASCADE,
		status BOOLEAN NOT NULL
	);
	`
	if _, err = DB.Exec(createTables); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}
