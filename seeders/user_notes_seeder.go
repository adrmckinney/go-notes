package seeders

import (
	"fmt"

	"github.com/adrmckinney/go-notes/db"
	"github.com/adrmckinney/go-notes/factories"
	"github.com/adrmckinney/go-notes/repos"
)

func SeedUserNotes(userCount, noteCount int) error {
	userIDs, err := SeedUsers(userCount)
	if err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}
	noteRepo := repos.NoteRepo{DB: db.GormDB}
	notes := factories.NoteFactory(noteCount, userIDs, "", "")
	for _, note := range notes {
		if _, err := noteRepo.CreateNote(note); err != nil {
			fmt.Printf("Failed to seed note: %v\n", err)
		}
	}
	return nil
}
