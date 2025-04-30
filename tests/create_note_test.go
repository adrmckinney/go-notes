package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adrmckinney/go-notes/models"
	"github.com/gorilla/mux"
)

func TestCreateNote(t *testing.T) {
	// Ensure the database is cleaned up after the test
	CleanUpDatabases(t, TestDB, []string{"notes"})

	body := models.Note{
		Title:   "test title",
		Content: "test content",
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to encode request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/notes", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/notes", NoteHandler.CreateNote).Methods("POST")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var createdNote models.Note
	err = json.Unmarshal(rr.Body.Bytes(), &createdNote)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Verify the note was inserted into the database
	_, err = NoteRepo.GetNoteById(uint(createdNote.ID))
	if err != nil {
		t.Fatalf("Failed to fetch note from database: %v", err)
	}

}
