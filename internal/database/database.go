package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// DB represents the database connection
type DB struct {
	conn *sql.DB
}

// NewDB creates a new database instance
func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Initialize the database schema
	if err := initSchema(db); err != nil {
		return nil, err
	}

	return &DB{conn: db}, nil
}

// initSchema initializes the database schema
func initSchema(db *sql.DB) error {
	// Create programs table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS programs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tid INTEGER UNIQUE NOT NULL,
			sid INTEGER DEFAULT 0,
			ch_id INTEGER DEFAULT 0,
			title TEXT NOT NULL,
			short_title TEXT,
			title_yomi TEXT,
			last_update TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	// Create program_cache table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS program_cache (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			service_id INTEGER NOT NULL,
			start_time DATETIME NOT NULL,
			duration INTEGER NOT NULL,
			tid INTEGER NOT NULL,
			title TEXT NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (tid) REFERENCES programs(tid),
			UNIQUE(service_id, start_time, duration)
		);
	`)
	if err != nil {
		return err
	}

	// Create indexes for better performance
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_programs_tid ON programs(tid);
		CREATE INDEX IF NOT EXISTS idx_program_cache_service ON program_cache(service_id, start_time, duration);
		CREATE INDEX IF NOT EXISTS idx_program_cache_expires ON program_cache(expires_at);
	`)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Conn returns the underlying database connection
func (db *DB) Conn() *sql.DB {
	return db.conn
}
