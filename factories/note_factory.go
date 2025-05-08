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
//   - count: the number of notes to generate
//   - userIDs: a slice of user IDs to assign to notes (rotated if count > len(userIDs))
//   - title, content: (ignored) placeholders for compatibility
//
// Returns:
//   - []models.Note: a slice of generated Note structs with random titles, content, and rotated user IDs
func NoteFactory(count int, userIDs []uint, title string, content string) []models.Note {
	notes := make([]models.Note, count)
	titles, sentences, err := utils.GetRandomSentences(count)
	if err != nil {
		fmt.Printf("Failed to fetch random sentences. Setting to empty string for default value. Error: %v\n", err)
		sentences = make([]string, count)
		titles = make([]string, count)
	}

	for i := range count {
		userID := uint(0)
		if len(userIDs) > 0 {
			userID = userIDs[i%len(userIDs)]
		}
		notes[i] = models.Note{
			UserID:  userID,
			Title:   titles[i],
			Content: sentences[i],
		}
	}

	return notes
}
