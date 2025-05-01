package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/repos"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type NoteHandler struct {
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
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	note, err := h.NoteRepo.GetNoteById(uint(id))
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
	createdNote, err := h.NoteRepo.CreateNote(note)

	if err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdNote)
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updateData map[string]any
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	filtered := models.FilterUpsertFields(updateData, models.AllowedNoteUpdateFields)

	if len(updateData) == 0 {
		http.Error(w, "No update data provided", http.StatusBadRequest)
		return
	}

	updatedNote, err := h.NoteRepo.UpdateNote(uint(id), filtered)

	if err != nil {
		// Check for GORM's record not found error
		if err == gorm.ErrRecordNotFound ||
			(err.Error() != "" && (err.Error() == "note not found" ||
				err.Error() == "note not found: record not found")) {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedNote)
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing 'id' path parameter", http.StatusBadRequest)
		return
	}

	err := h.NoteRepo.DeleteNote(id)

	if err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Note successfully deleted",
	})
}
