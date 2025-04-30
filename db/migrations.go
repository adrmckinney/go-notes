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
	// Use the existing GORM DB connection
	sqlDB, err := GormDB.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from GormDB: %v", err)
	}

	// Configure the migration driver
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
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

	// Apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")
}

func RunRollback(steps int) {
	sqlDB, err := GormDB.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from GormDB: %v", err)
	}

	// Configure the migration driver
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
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
		log.Println("Rolling back all migrations...")
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback all migrations: %v", err)
		}
		log.Println("All migrations rolled back successfully!")
	} else {
		log.Printf("Rolling back %d migration steps...", steps)
		for i := 0; i < steps; i++ {
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

	// Remove the database using a server-level connection
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

func EnsureDatabaseExists() {
	cfg := config.GetDBConfig()
	serverDSN := config.GetServerDSN() // DSN without database name
	dbName := cfg.Database

	serverDB, err := sql.Open("mysql", serverDSN)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL server: %v", err)
	}
	defer serverDB.Close()

	_, err = serverDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		log.Fatalf("Failed to create database %s: %v", dbName, err)
	}
}
