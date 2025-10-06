package storage

import "database/sql"

// SQLite is a lightweight holder for the shared database handle.
type SQLite struct {
	DB *sql.DB
}

// NewSQLite wires the raw handle into the storage layer; real migrations arrive later.
func NewSQLite(db *sql.DB) *SQLite {
	return &SQLite{DB: db}
}
