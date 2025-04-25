package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adrmckinney/go-notes/factories"
	"github.com/adrmckinney/go-notes/models"
	"github.com/gorilla/mux"
)

func TestGetNote(t *testing.T) {
	// Ensure the database is cleaned up after the test
	CleanUpDatabases(t, TestDB, []string{"notes"})

	notes := factories.NoteFactory(2, "", "")
	expectedNote := notes[1]

	// Seed the database with test data
	for _, note := range notes {
		_, err := NoteRepo.CreateNote(note)
		if err != nil {
			t.Fatalf("Failed to seed test database: %v", err)
		}
	}

	// Create a new mux.Router and register the route so that the path param will be readable
	router := mux.NewRouter()
	router.HandleFunc("/notes/{id}", NoteHandler.GetNote).Methods("GET")

	t.Run("Successful fetch", func(t *testing.T) {
		// Create a new HTTP request
		req, err := http.NewRequest("GET", "/notes/2", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		router.ServeHTTP(rr, req)

		// Check the response status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Decode the response body
		var resNote models.Note
		err = json.Unmarshal(rr.Body.Bytes(), &resNote)
		if err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		if resNote.Title != expectedNote.Title {
			t.Fatalf("Failed to get correct note: expected title %s, got title %s", expectedNote.Title, resNote.Title)
		}
	})

	t.Run("Note not found", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/notes/999", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

}
