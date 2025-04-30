package routes

import (
	"github.com/adrmckinney/go-notes/handlers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()

	// Implemented init handlers because testing requires a
	// sqlite db to be created and used, which means we need to
	// pass in the correct DB here.
	handlers := handlers.InitHandlers(db)

	r.HandleFunc("/notes", handlers.NoteHandler.GetNotes).Methods("GET")
	r.HandleFunc("/notes/{id}", handlers.NoteHandler.GetNote).Methods("GET")
	r.HandleFunc("/notes", handlers.NoteHandler.CreateNote).Methods("POST")
	r.HandleFunc("/notes/{id}", handlers.NoteHandler.UpdateNote).Methods("PUT")
	r.HandleFunc("/notes/{id}", handlers.NoteHandler.DeleteNote).Methods("DELETE")
	return r
}
