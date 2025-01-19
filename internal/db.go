package internal

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const DefaultDBName = "iknow.sqlite3"

func openDB(dbName string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to open the database: %w", err)
	}
	return db, nil
}

func initDB(db *sqlx.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS entries (date CHAR(16) PRIMARY KEY, started_items INTEGER, completed_items INTEGER, completed_courses INTEGER, study_time INTEGER)`,
	}
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute the DB init query: %w", err)
		}
	}
	return nil
}
