package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/adrmckinney/go-notes/factories"
	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/routes"
)

func TestGetNotes(t *testing.T) {
	TearDown(t)

	user := InitUser(t, InitUserOptions{})
	otherUser := InitUser(t, InitUserOptions{})
	notes := factories.NoteFactory(2, []uint{user.ID}, "", "")
	factories.NoteFactory(2, []uint{otherUser.ID}, "", "")

	// Seed the database with test data
	for _, note := range notes {
		_, err := NoteRepo.CreateNote(note)
		if err != nil {
			t.Fatalf("Failed to seed test database: %v", err)
		}
	}

	rr := CreateRouteAndServe(t, routes.GET_NOTES, ServeOpts{AuthToken: &user.Token})

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Decode the response body
	var resNotes []models.Note
	err := json.Unmarshal(rr.Body.Bytes(), &resNotes)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Compare the length of the original notes and the response notes
	if len(resNotes) != len(notes) {
		t.Fatalf("Length mismatch: expected %d, got %d", len(notes), len(resNotes))
	}

	allMatch := true
	for _, note := range resNotes {
		if note.UserID != user.ID {
			allMatch = false
			break
		}

		if !allMatch {
			t.Fatalf("Get notes query returned notes belonging to other users")
		}
	}
}
