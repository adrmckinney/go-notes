package handlers

import (
	"github.com/adrmckinney/go-notes/repos"
	"gorm.io/gorm"
)

type Handlers struct {
	NoteHandler *NoteHandler
}

func InitHandlers(db *gorm.DB) *Handlers {
	noteRepo := &repos.NoteRepo{DB: db}

	// Initialize handlers with repositories
	noteHandler := &NoteHandler{
		NoteRepo: noteRepo,
	}

	return &Handlers{
		NoteHandler: noteHandler,
	}
}
