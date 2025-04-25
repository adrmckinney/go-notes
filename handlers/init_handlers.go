package handlers

import (
	"database/sql"
	"log"

	"github.com/adrmckinney/go-notes/repos"
)

type Handlers struct {
	NoteHandler *NoteHandler
}

func InitHandlers(db *sql.DB) *Handlers {
	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection is invalid: %v", err)
	}

	noteRepo := &repos.NoteRepo{DB: db}

	// Initialize handlers with repositories
	noteHandler := &NoteHandler{
		NoteRepo: noteRepo,
	}

	return &Handlers{
		NoteHandler: noteHandler,
	}
}
