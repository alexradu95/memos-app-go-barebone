package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

var (
	Db      *sql.DB
	once    sync.Once
	initErr error
)

func Initialize() error {

	once.Do(func() {
		primaryURL := os.Getenv("TURSO_DATABASE_URL")
		authToken := os.Getenv("TURSO_AUTH_TOKEN")

		var connString string
		if primaryURL != "" && authToken != "" {
			// Use Turso database if credentials are available
			connString = fmt.Sprintf("%s?authToken=%s", primaryURL, authToken)
		} else {
			// Fall back to local SQLite database
			log.Println("Turso credentials not found, using local SQLite database")
			connString = "file:local.db"
		}

		var db *sql.DB
		db, initErr = sql.Open("libsql", connString)
		if initErr != nil {
			initErr = fmt.Errorf("failed to open db (%s): %w", primaryURL, initErr)
			log.Println(initErr)
			return
		}

		if pingErr := db.Ping(); pingErr != nil {
			initErr = fmt.Errorf("failed to ping database: %w", pingErr)
			log.Println(initErr)
			return
		}

		// If using local database, ensure it exists and has necessary tables
		if primaryURL == "" {
			if err := initializeLocalDB(db); err != nil {
				initErr = fmt.Errorf("failed to initialize local database: %w", err)
				log.Println(initErr)
				return
			}
		}

		Db = db
	})

	return initErr
}

// initializeLocalDB creates necessary tables in the local database
func initializeLocalDB(db *sql.DB) error {
	// Add your table creation statements here
	// Example:
	_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS accounts (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		username TEXT NOT NULL,
    		password_hash TEXT NOT NULL,
    		created_at TEXT NOT NULL);`)

	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS posts (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		content TEXT,
    		created_at TEXT NOT NULL,
    		updated_at TEXT NOT NULL,
    		account_id INTEGER NOT NULL,
    		FOREIGN KEY (account_id) REFERENCES accounts(id));`)

	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
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
