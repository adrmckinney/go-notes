package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/adrmckinney/go-notes/auth"
	"github.com/adrmckinney/go-notes/errs"
	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/services"
	"github.com/gorilla/mux"
)

type NoteHandler struct {
	NoteService services.NoteService
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	authUser, ok := auth.GetAuthInfo(r)
	if !ok {
		http.Error(w, "Unable to find authUser", http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	result, err := h.NoteService.GetNoteById(authUser.UserID, uint(id))
	if err != nil {
		if httpErr, ok := err.(*errs.HTTPError); ok {
			http.Error(w, httpErr.Message, httpErr.Status)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (h *NoteHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	authUser, ok := auth.GetAuthInfo(r)
	if !ok {
		http.Error(w, "Unable to find authUser", http.StatusInternalServerError)
		return
	}
	notes, err := h.NoteService.GetNotes(authUser.UserID)
	if err != nil {
		if httpErr, ok := err.(*errs.HTTPError); ok {
			http.Error(w, httpErr.Message, httpErr.Status)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	authInfo, ok := auth.GetAuthInfo(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload models.Note
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdNote, err := h.NoteService.CreateNote(authInfo.UserID, payload)
	if err != nil {
		if httpErr, ok := err.(*errs.HTTPError); ok {
			http.Error(w, httpErr.Message, httpErr.Status)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdNote)
}

// HERE

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	authInfo, ok := auth.GetAuthInfo(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	noteId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updateData map[string]any
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	filtered := models.FilterUpdateFields(updateData, models.AllowedNoteUpdateFields)

	if len(updateData) == 0 {
		http.Error(w, "No update data provided", http.StatusBadRequest)
		return
	}

	result, err := h.NoteService.UpdateNote(authInfo.UserID, uint(noteId), filtered)
	if err != nil {
		if httpErr, ok := err.(*errs.HTTPError); ok {
			http.Error(w, httpErr.Message, httpErr.Status)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	authInfo, ok := auth.GetAuthInfo(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	noteId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := h.NoteService.DeleteNote(authInfo.UserID, uint(noteId))
	if err != nil {
		if httpErr, ok := err.(*errs.HTTPError); ok {
			http.Error(w, httpErr.Message, httpErr.Status)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

}
