package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "modernc.org/sqlite" // Import the SQLite driver
)

var (
	Db    *sql.DB
	once  sync.Once
	errDB error // Renamed to avoid shadowing the 'err' inside once.Do
)

func Initialize() error {
	once.Do(func() {
		// Use a local SQLite database file.
		connString := "file:local.db?_pragma=foreign_keys(1)" // Enable foreign keys

		var db *sql.DB
		db, errDB = sql.Open("sqlite", connString)
		if errDB != nil {
			errDB = fmt.Errorf("failed to open db (local.db): %w", errDB)
			log.Println(errDB)
			return // Return from the anonymous function, setting errDB
		}

		if errDB = db.Ping(); errDB != nil {
			errDB = fmt.Errorf("failed to ping database: %w", errDB)
			log.Println(errDB)
			return
		}

		// Initialize the database schema (create tables if they don't exist)
		if errDB = initializeLocalDB(db); errDB != nil {
			errDB = fmt.Errorf("failed to initialize local database: %w", errDB)
			log.Println(errDB)
			return
		}

		Db = db
	})

	return errDB
}

// initializeLocalDB creates necessary tables in the local database
func initializeLocalDB(db *sql.DB) error {
	// Create the 'accounts' table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,  -- Added UNIQUE constraint
			password_hash TEXT NOT NULL,
			created_at TEXT NOT NULL
		);`)
	if err != nil {
		return fmt.Errorf("failed to create accounts table: %w", err)
	}

	// Create the 'posts' table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			account_id INTEGER NOT NULL,
			FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE -- Added ON DELETE CASCADE
		);`)
	if err != nil {
		return fmt.Errorf("failed to create posts table: %w", err)
	}

	return nil
}

// CloseDB closes the global database connection if open.
func CloseDB() error {
	if Db != nil {
		return Db.Close()
	}
	return nil
}
