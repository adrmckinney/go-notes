package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/adrmckinney/go-notes/factories"
	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/routes"
)

func TestUpdateNote(t *testing.T) {
	TearDown(t)

	user := InitUser(t, InitUserOptions{})
	notes := factories.NoteFactory(1, []uint{1}, "", "")
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

		rr := CreateRouteAndServe(t, routes.UPDATE_NOTE, ServeOpts{AuthToken: &user.Token, Payload: body, PathParams: map[string]any{"id": 1}})

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

		rr := CreateRouteAndServe(t, routes.UPDATE_NOTE, ServeOpts{AuthToken: &user.Token, Payload: body, PathParams: map[string]any{"id": 999}})

		// Check the response status code
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}

	})
}
