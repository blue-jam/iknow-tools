package internal

import (
	"fmt"
	"github.com/urfave/cli/v2"

	"github.com/jmoiron/sqlx"
)

const (
	DefaultDBName = "iknow.sqlite3"
	dbMetaKey     = "db"
)

func OpenDBFromContext(ctx *cli.Context) error {
	dbName := ctx.String(dbNameFlag.Name)
	if dbName == "" {
		dbName = DefaultDBName
	}
	db, err := openDB(dbName)
	if err != nil {
		return err
	}

	err = initDB(db)
	if err != nil {
		return err
	}

	SetDBToContext(ctx, db)

	return nil
}

func SetDBToContext(ctx *cli.Context, db *sqlx.DB) {
	ctx.App.Metadata[dbMetaKey] = db
}

func GetDBFromContext(ctx *cli.Context) *sqlx.DB {
	return ctx.App.Metadata[dbMetaKey].(*sqlx.DB)
}

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
