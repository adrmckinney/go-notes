package factories

import (
	"fmt"

	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/utils"
)

// NoteFactory generates a slice of models.Note for testing or seeding purposes.
// It fetches random titles and content using utils.GetRandomSentences. If fetching fails,
// it fills the content with empty strings. The function ignores the provided title and content
// arguments and instead uses the random data for each note.
//
// Params:
//
//	count   - the number of notes to generate
//	title   - (ignored) placeholder for compatibility
//	content - (ignored) placeholder for compatibility
//
// Returns:
//
//	[]models.Note - a slice of generated Note structs with random titles and content
func NoteFactory(count int, title string, content string) []models.Note {
	notes := make([]models.Note, count)
	titles, sentences, err := utils.GetRandomSentences(count)
	if err != nil {
		fmt.Printf("Failed to fetch random sentences. Setting to empty string for default value. Error: %v\n", err)
		sentences = make([]string, count)
	}

	for i := range count {
		notes[i] = models.Note{
			Title:   titles[i],
			Content: sentences[i],
		}
	}

	return notes
}
