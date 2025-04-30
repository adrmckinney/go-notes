package seeders

import (
	"fmt"

	"github.com/adrmckinney/go-notes/db"
	"github.com/adrmckinney/go-notes/factories"
	"github.com/adrmckinney/go-notes/repos"
)

func SeedNotes(count int) {
	noteRepo := repos.NoteRepo{DB: db.GormDB}

	notes := factories.NoteFactory(count, "", "")

	for _, note := range notes {
		_, err := noteRepo.CreateNote(note)
		if err != nil {
			fmt.Printf("Failed to seed note: %v\n", err)
		}
	}
}
