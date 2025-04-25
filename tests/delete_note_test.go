package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adrmckinney/go-notes/factories"
	"github.com/gorilla/mux"
)

func TestDeleteNote(t *testing.T) {
	// Ensure the database is cleaned up after the test
	CleanUpDatabases(t, TestDB, []string{"notes"})

	notes := factories.NoteFactory(1, "", "")
	note := notes[0]
	_, err := NoteRepo.CreateNote(note)
	if err != nil {
		t.Fatalf("Failed to seed test database: %v", err)
	}

	// Create a new mux.Router and register the route so that the path param will be readable
	router := mux.NewRouter()
	router.HandleFunc("/notes/{id}", NoteHandler.DeleteNote).Methods("DELETE")

	// Create a new HTTP request
	req, err := http.NewRequest("DELETE", "/notes/1", nil)
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

	// Verify the note was deleted from the database
	_, err = NoteRepo.GetNoteById("1")
	if err == nil {
		t.Fatalf("Note was not deleted from the database")
	}
}
