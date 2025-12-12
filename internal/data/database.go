package data

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// InitDB connects to the database and returns a new connection.
func InitDB() (*sqlx.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "internal/data/mainframe.db"
	}
	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Enable foreign keys for SQLite
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	return db, nil
}
