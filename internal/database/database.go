package database

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/libsql/libsql-client-go"
)

var (
	Db   *sql.DB
	once sync.Once
	err  error
)

func Initialize() error {
	once.Do(func() {
		// Get Turso connection details from environment variables
		dbURL := os.Getenv("TURSO_DATABASE_URL")
		authToken := os.Getenv("TURSO_AUTH_TOKEN")

		// If Turso environment variables are not set, fall back to SQLite
		if dbURL == "" {
			Db, err = sql.Open("sqlite3", "./data.db")
			if err != nil {
				log.Printf("Failed to connect to SQLite database: %v", err)
				return
			}
		} else {
			// Connect to Turso
			connectionString := dbURL
			if authToken != "" {
				connectionString += "?authToken=" + authToken
			}

			Db, err = sql.Open("libsql", connectionString)
			if err != nil {
				log.Printf("Failed to connect to Turso database: %v", err)
				return
			}
		}

		if err = Db.Ping(); err != nil {
			log.Printf("Failed to ping database: %v", err)
		}
	})
	return err
}

func CloseDB() error {
	if Db != nil {
		return Db.Close()
	}
	return nil
}
