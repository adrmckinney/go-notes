package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adrmckinney/go-notes/factories"
	"github.com/adrmckinney/go-notes/models"
	"github.com/gorilla/mux"
)

func TestUpdateNote(t *testing.T) {
	// Ensure the database is cleaned up after the test
	CleanUpDatabases(t, TestDB, []string{"notes"})

	notes := factories.NoteFactory(1, "", "")
	note := notes[0]
	_, err := NoteRepo.CreateNote(note)
	if err != nil {
		t.Fatalf("Failed to seed test database: %v", err)
	}

	t.Run("Successful title update", func(t *testing.T) {
		newTitle := "updated title"
		body := models.Note{
			Title: newTitle,
		}
		reqBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to encode request body: %v", err)
		}

		req, err := http.NewRequest("PUT", "/notes/1", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		router := mux.NewRouter()
		router.HandleFunc("/notes/{id}", NoteHandler.UpdateNote).Methods("PUT")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var updatedNote models.Note
		err = json.Unmarshal(rr.Body.Bytes(), &updatedNote)
		if err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		if updatedNote.Title != newTitle {
			t.Fatalf("Failed to update title: Expected %s Got %s", newTitle, updatedNote.Title)
		}
	})

	t.Run("Error missing update data", func(t *testing.T) {

		body := models.Note{
			Content: "",
		}
		reqBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to encode request body: %v", err)
		}

		// Create a new HTTP request
		req, err := http.NewRequest("PUT", "/notes/999", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		router := mux.NewRouter()
		router.HandleFunc("/notes/{id}", NoteHandler.UpdateNote).Methods("PUT")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the response status code
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}

	})
}
