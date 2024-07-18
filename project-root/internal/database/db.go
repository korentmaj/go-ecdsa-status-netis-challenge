package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	return nil
}
