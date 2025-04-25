package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/repos"
	"github.com/gorilla/mux"
)

type NoteHandler struct {
	// DB *sql.DB
	NoteRepo *repos.NoteRepo
}

// Go's http.Request holds a lot of data. It would be inefficient
// to copy the Request. Passing the reference is more efficient.
func (h *NoteHandler) GetNotes(w http.ResponseWriter, _ *http.Request) {
	notes, err := h.NoteRepo.GetNotes()
	if err != nil {
		http.Error(w, "Failed to fetch notes", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing 'id' path parameter", http.StatusBadRequest)
		return
	}
	note, err := h.NoteRepo.GetNoteById(id)
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
	}
	json.NewEncoder(w).Encode(note)
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil { // get an explanation for this syntax
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if note.Title == "" || note.Content == "" {
		http.Error(w, "Title and Content are required", http.StatusBadRequest)
		return
	}
	id, err := h.NoteRepo.CreateNote(note)

	if err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	createdNote, err := h.NoteRepo.GetNoteById(fmt.Sprintf("%d", id))
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdNote)
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing 'id' path parameter", http.StatusBadRequest)
		return
	}

	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if note.Title == "" && note.Content == "" {
		http.Error(w, "At least one of 'title' or 'content' must be provided", http.StatusBadRequest)
		return
	}

	err := h.NoteRepo.UpdateNote(id, note)
	if err != nil {
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
	}

	updatedNote, err := h.NoteRepo.GetNoteById(id)
	if err != nil {
		http.Error(w, "Failed to retrieve updated note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedNote)
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing 'id' path parameter", http.StatusBadRequest)
		return
	}

	message, err := h.NoteRepo.DeleteNote(id)

	if err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}
