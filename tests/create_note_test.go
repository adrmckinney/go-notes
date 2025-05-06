package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/routes"
)

func TestCreateNote(t *testing.T) {
	TearDown(t)

	user := InitUser(t, InitUserOptions{})

	body := models.Note{
		Title:   "test title",
		Content: "test content",
	}

	rr := CreateRouteAndServe(t, routes.CREATE_NOTE, ServeOpts{Payload: body, AuthToken: &user.Token})

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdNote models.Note
	err := json.Unmarshal(rr.Body.Bytes(), &createdNote)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Verify the note was inserted into the database
	_, err = NoteRepo.GetNoteById(uint(createdNote.ID))
	if err != nil {
		t.Fatalf("Failed to fetch note from database: %v", err)
	}
}
