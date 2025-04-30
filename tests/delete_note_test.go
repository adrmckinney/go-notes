package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adrmckinney/go-notes/factories"
	"github.com/gorilla/mux"
)

func TestDeleteNote(t *testing.T) {
	CleanUpDatabases(t, TestDB, []string{"notes"})

	notes := factories.NoteFactory(1, "", "")
	note := notes[0]
	_, err := NoteRepo.CreateNote(note)
	if err != nil {
		t.Fatalf("Failed to seed test database: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/notes/{id}", NoteHandler.DeleteNote).Methods("DELETE")

	req, err := http.NewRequest("DELETE", "/notes/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	_, err = NoteRepo.GetNoteById(1)
	if err == nil {
		t.Fatalf("Note was not deleted from the database")
	}
}
