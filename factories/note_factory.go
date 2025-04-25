package factories

import (
	"fmt"
	"time"

	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/utils"
)

func NoteFactory(count int, title string, content string) []models.Note {
	notes := make([]models.Note, count)
	titles, sentences, err := utils.GetRandomSentences(count)
	if err != nil {
		fmt.Printf("Failed to fetch random sentences. Setting to empty string for default value. Error: %v\n", err)
		sentences = make([]string, count)
	}

	for i := range count {
		now := time.Now().UTC()
		notes[i] = models.Note{
			Title:    titles[i],
			Content:  sentences[i],
			Added:    now.Format(time.RFC3339),
			Modified: now.Format(time.RFC3339),
		}
	}

	return notes
}
