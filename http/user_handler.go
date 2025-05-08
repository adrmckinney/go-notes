package http

import (
	"encoding/json"
	"net/http"

	"github.com/adrmckinney/go-notes/errs"
	"github.com/adrmckinney/go-notes/services"
)

type UserHandler struct {
	UserService services.UserService
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := GetIdFromUrlPath(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := h.UserService.GetUserById(id)
	if err != nil {
		if httpErr, ok := err.(*errs.HTTPError); ok {
			http.Error(w, httpErr.Message, httpErr.Status)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := GetIdFromUrlPath(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var updateData map[string]any
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(updateData) == 0 {
		http.Error(w, "No update data provided", http.StatusBadRequest)
		return
	}

	result, err := h.UserService.UpdateUser(id, updateData)
	if err != nil {
		if httpErr, ok := err.(*errs.HTTPError); ok {
			http.Error(w, httpErr.Message, httpErr.Status)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := GetIdFromUrlPath(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := h.UserService.DeleteUser(id)
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
