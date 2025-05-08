package seeders

import (
	"fmt"

	"github.com/adrmckinney/go-notes/factories"
)

func SeedUserNotes(deps SeederDeps, userCount, noteCount int) error {
	userIDs, err := SeedUsers(deps, userCount)
	if err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	notes := factories.NoteFactory(noteCount, userIDs, "", "")
	for i, note := range notes {
		targetIdx := i % len(userIDs)
		if _, err := deps.NoteService.CreateNote(userIDs[targetIdx], note); err != nil {
			fmt.Printf("Failed to seed note: %v\n", err)
		}
	}
	return nil
}
