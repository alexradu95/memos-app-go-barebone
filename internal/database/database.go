package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Db   *sql.DB
	once sync.Once
	err  error
)

func Initialize() error {
	once.Do(func() {
		Db, err = sql.Open("sqlite3", "./data.db")
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			return
		}
		// Verify the connection is successful
		if err = Db.Ping(); err != nil {
			log.Printf("Failed to ping database: %v", err)
		}
	})
	return err
}

// CloseDB closes the database connection.
func CloseDB() error {
	if Db != nil {
		return Db.Close()
	}
	return nil
}
