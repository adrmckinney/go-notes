package tests

import (
	"database/sql"
	"fmt"
	"testing"

	"gorm.io/gorm"
)

func CleanUpTables(t *testing.T, testDB *gorm.DB, tables []string) {
	for _, table := range tables {
		tbl := table // capture loop variable
		t.Cleanup(func() {
			query := fmt.Sprintf("DELETE FROM %s", tbl)
			if err := testDB.Exec(query).Error; err != nil {
				t.Errorf("Failed to clean up table %s: %v", tbl, err)
			}

			// Reset AUTOINCREMENT counter for SQLite
			if err := testDB.Exec(fmt.Sprintf("DELETE FROM sqlite_sequence WHERE name='%s'", tbl)).Error; err != nil {
				t.Errorf("Failed to reset AUTOINCREMENT for table %s: %v", tbl, err)
			}
		})
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
