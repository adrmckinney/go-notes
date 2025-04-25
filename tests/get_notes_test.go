package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adrmckinney/go-notes/factories"
	"github.com/adrmckinney/go-notes/models"
)

func TestGetNotes(t *testing.T) {
	// Ensure the database is cleaned up after the test
	CleanUpDatabases(t, TestDB, []string{"notes"})

	notes := factories.NoteFactory(2, "", "")

	// Seed the database with test data
	for _, note := range notes {
		_, err := NoteRepo.CreateNote(note)
		if err != nil {
			t.Fatalf("Failed to seed test database: %v", err)
		}
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/notes", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler

	NoteHandler.GetNotes(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Decode the response body
	var resNotes []models.Note
	err = json.Unmarshal(rr.Body.Bytes(), &resNotes)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Compare the length of the original notes and the response notes
	if len(resNotes) != len(notes) {
		t.Fatalf("Length mismatch: expected %d, got %d", len(notes), len(resNotes))
	}
}
