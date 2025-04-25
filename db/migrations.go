package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/adrmckinney/go-notes/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
)

func RunMigrations() {
	// Connect to the MySQL server
	serverDB, err := sql.Open("mysql", config.GetServerDSN())
	if err != nil {
		log.Fatalf("Failed to connect to MySQL server: %v", err)
	}
	defer serverDB.Close()

	// Create the database if it doesn't exist
	dbName := config.GetDBConfig().Database
	_, err = serverDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		log.Fatalf("Failed to create database %s: %v", dbName, err)
	}
	log.Printf("Database %s is ready.", dbName)

	// Connect to the specific database
	db, err := sql.Open("mysql", config.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database %s: %v", dbName, err)
	}
	defer db.Close()

	// Configure the migration driver
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	// Run migrations from the "migrations" directory
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql", driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	// Check the current migration version
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Failed to check migration version: %v", err)
	}

	if dirty {
		log.Fatalf("Database is in a dirty state at version %d. Please resolve the issue before proceeding.", version)
	}

	if err == nil {
		log.Printf("Migrations already applied. Current version: %d", version)
		return
	}

	// Apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")
}

func RunRollback(steps int) {
	serverDB, err := sql.Open("mysql", config.GetServerDSN())
	if err != nil {
		log.Fatalf("Failed to connect to MySQL server: %v", err)
	}
	defer serverDB.Close()

	db, err := sql.Open("mysql", config.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database %s: %v", config.GetDBConfig().Database, err)
	}
	defer db.Close()

	// Configure the migration driver
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	// Run migrations from the "migrations" directory
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql", driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	// Rollback logic
	if steps == 0 {
		// Roll back all migrations
		log.Println("Rolling back all migrations...")
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback all migrations: %v", err)
		}
		log.Println("All migrations rolled back successfully!")
	} else {
		// Roll back the specified number of steps
		log.Printf("Rolling back %d migration steps...", steps)
		for i := range steps {
			if err := m.Steps(-1); err != nil {
				if err == migrate.ErrNoChange {
					log.Println("No more migrations to rollback.")
					break
				}
				log.Fatalf("Failed to rollback migrations: %v", err)
			}
			log.Printf("Rolled back 1 step. Remaining steps: %d", steps-i-1)
		}
		log.Println("Rollback completed successfully!")
	}
}

func RunRemoveDatabase() {
	RunRollback(0)

	serverDB, err := sql.Open("mysql", config.GetServerDSN())
	if err != nil {
		log.Fatalf("Failed to connect to MySQL server: %v", err)
	}
	defer serverDB.Close()

	dbName := config.GetDBConfig().Database
	_, err = serverDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	if err != nil {
		log.Fatalf("Failed to drop database %s: %v", dbName, err)
	}

	log.Printf("Database %s has been removed successfully.", dbName)
}
