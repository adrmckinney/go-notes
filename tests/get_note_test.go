package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/adrmckinney/go-notes/factories"
	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/routes"
)

func TestGetNote(t *testing.T) {
	TearDown(t)

	notes := factories.NoteFactory(2, []uint{1}, "", "")
	expectedNote := notes[1]

	user := InitUser(t, InitUserOptions{})

	// Seed the database with test data
	for _, note := range notes {
		_, err := NoteRepo.CreateNote(note)
		if err != nil {
			t.Fatalf("Failed to seed test database: %v", err)
		}
	}

	t.Run("Successful fetch", func(t *testing.T) {
		rr := CreateRouteAndServe(t, routes.GET_NOTE, ServeOpts{PathParams: map[string]any{"id": 2}, AuthToken: &user.Token})

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var resNote models.Note
		err := json.Unmarshal(rr.Body.Bytes(), &resNote)
		if err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		if resNote.Title != expectedNote.Title {
			t.Fatalf("Failed to get correct note: expected title %s, got title %s", expectedNote.Title, resNote.Title)
		}
	})

	t.Run("Note not found", func(t *testing.T) {
		rr := CreateRouteAndServe(t, routes.GET_NOTE, ServeOpts{PathParams: map[string]any{"id": 999}, AuthToken: &user.Token})

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

}
