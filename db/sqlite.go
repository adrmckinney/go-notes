package db

import (
	"fmt"
	"log"
	"testing"

	"github.com/adrmckinney/go-notes/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func InitTestGorm() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to sqlite: %v", err)
	}
	db.AutoMigrate(
		&models.User{},
		&models.Note{},
		&models.UserToken{},
	)
	return db
}

func CleanUpAllTables(t *testing.T, testDB *gorm.DB) {
	t.Helper()

	t.Cleanup(func() {
		var tables []string

		// Query all user-defined tables (exclude sqlite_ system tables)
		err := testDB.Raw(`
			SELECT name FROM sqlite_master
			WHERE type='table' AND name NOT LIKE 'sqlite_%';
		`).Scan(&tables).Error

		if err != nil {
			t.Errorf("Failed to fetch table names for cleanup: %v", err)
			return
		}

		// Disable FK constraints temporarily to avoid constraint issues during cleanup
		testDB.Exec("PRAGMA foreign_keys = OFF;")

		for _, table := range tables {
			if err := testDB.Exec(fmt.Sprintf("DELETE FROM %s;", table)).Error; err != nil {
				t.Errorf("Failed to delete records from table %s: %v", table, err)
			}
			if err := testDB.Exec(fmt.Sprintf("DELETE FROM sqlite_sequence WHERE name='%s';", table)).Error; err != nil {
				t.Errorf("Failed to reset AUTOINCREMENT for table %s: %v", table, err)
			}
		}

		testDB.Exec("PRAGMA foreign_keys = ON;")
	})
}
