package tests

import (
	"database/sql"
	"fmt"
	"testing"

	"gorm.io/gorm"
)

func CleanUpDatabases(t *testing.T, testDB *gorm.DB, tables []string) {
	for _, table := range tables {
		t.Cleanup(func() {
			query := fmt.Sprintf("DELETE FROM %s", table)
			if err := testDB.Exec(query).Error; err != nil {
				t.Fatalf("Failed to clean up test database: %v", err)
			}
		})

		// Reset AUTOINCREMENT counter for SQLite
		if err := testDB.Exec(fmt.Sprintf("DELETE FROM sqlite_sequence WHERE name='%s'", table)).Error; err != nil {
			t.Fatalf("Failed to reset AUTOINCREMENT counter for table %s: %v", table, err)
		}
	}
}

func LogTestDatabaseState(t *testing.T, testDB *sql.DB) {
	rows, err := testDB.Query("SELECT id, title, content FROM notes")
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	defer rows.Close()

	t.Log("Current database state:")
	for rows.Next() {
		var id int
		var title, content string
		if err := rows.Scan(&id, &title, &content); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		t.Logf("ID: %d, Title: %s, Content: %s", id, title, content)
	}
}
