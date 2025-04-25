package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func InitTestDB() *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("Failed to create test database: %v", err)
	}

	// Run migrations
	err = runMigrations(db)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func runMigrations(db *sql.DB) error {
	// Example schema creation
	queries := []string{
		`CREATE TABLE notes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            added DATETIME DEFAULT CURRENT_TIMESTAMP,
            modified DATETIME DEFAULT CURRENT_TIMESTAMP
        );`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	fmt.Println("Migrations completed successfully")
	return nil
}
