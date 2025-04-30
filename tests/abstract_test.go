package tests

import (
	"os"
	"testing"

	"github.com/adrmckinney/go-notes/db"
	"github.com/adrmckinney/go-notes/handlers"
	"github.com/adrmckinney/go-notes/repos"
	"gorm.io/gorm"
)

var TestDB *gorm.DB          // Shared test database connection
var NoteRepo *repos.NoteRepo // Shared NoteRepo instance
var NoteHandler *handlers.NoteHandler

// TestMain is the entry point for all tests in the tests package and its subpackages
func TestMain(m *testing.M) {
	// Initialize the test database
	TestDB = db.InitTestGorm()

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
