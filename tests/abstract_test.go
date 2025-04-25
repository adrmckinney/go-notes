package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/adrmckinney/go-notes/db"
	"github.com/adrmckinney/go-notes/handlers"
	"github.com/adrmckinney/go-notes/repos"
)

var TestDB *sql.DB           // Shared test database connection
var NoteRepo *repos.NoteRepo // Shared NoteRepo instance
var NoteHandler *handlers.NoteHandler

// TestMain is the entry point for all tests in the tests package and its subpackages
func TestMain(m *testing.M) {
	// Initialize the test database
	TestDB = db.InitTestDB()
	defer TestDB.Close()

	// Verify the database connection
	if err := TestDB.Ping(); err != nil {
		log.Fatalf("Test database connection is invalid: %v", err)
	}

	// Initialize the Repos
	// NoteRepo = &repos.NoteRepo{DB: TestDB}

	// // Initialize Handlers
	// NoteHandler = &handlers.NoteHandler{DB: TestDB}

	// Initialize repositories and handlers
	initializeTestDependencies()

	// Run the tests
	code := m.Run()

	// Exit with the test result code
	os.Exit(code)
}

// initializeTestDependencies initializes shared repositories and handlers for tests
func initializeTestDependencies() {
	// Initialize the NoteRepo
	NoteRepo = &repos.NoteRepo{DB: TestDB}

	// Initialize the NoteHandler with the NoteRepo
	NoteHandler = &handlers.NoteHandler{NoteRepo: NoteRepo}
}
