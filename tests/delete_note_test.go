package tests

import (
	"net/http"
	"testing"

	"github.com/adrmckinney/go-notes/factories"
	"github.com/adrmckinney/go-notes/routes"
)

func TestDeleteNote(t *testing.T) {
	TearDown(t)

	notes := factories.NoteFactory(1, []uint{1}, "", "")
	note := notes[0]
	_, err := NoteRepo.CreateNote(note)
	if err != nil {
		t.Fatalf("Failed to seed test database: %v", err)
	}
	user := InitUser(t, InitUserOptions{})
	rr := CreateRouteAndServe(
		t,
		routes.DELETE_NOTE,
		ServeOpts{
			AuthToken:  &user.Token,
			PathParams: map[string]any{"id": user.ID},
		})

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	_, err = NoteRepo.GetNoteById(1)
	if err == nil {
		t.Fatalf("Note was not deleted from the database")
	}
}
