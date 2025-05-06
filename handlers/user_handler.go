package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/repos"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserRepo *repos.UserRepo
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := h.UserRepo.GetUsers()
	if err != nil {
		http.Error(w, "failed to fetch users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := GetIdFromUrlPath(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	user, err := h.UserRepo.GetUserById(uint(id))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
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

	password, passOk := updateData["password"].(string)
	confirmPassword, confirmOk := updateData["confirmPassword"].(string)

	if passOk && (!confirmOk || confirmPassword == "") {
		http.Error(w, "ConfirmPassword must be provided for password change", http.StatusBadRequest)
		return
	}

	if passOk && confirmOk && password != confirmPassword {
		http.Error(w, "Password and ConfirmPassword do not match", http.StatusBadRequest)
		return
	}

	filtered := models.FilterUpdateFields(updateData, models.AllowedUserUpdateFields)

	if passRaw, ok := filtered["password"]; ok && passRaw != "" {
		if passStr, ok := passRaw.(string); ok {
			hashedPass, err := bcrypt.GenerateFromPassword([]byte(passStr), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Failed to hash password", http.StatusInternalServerError)
				return
			}
			filtered["password"] = string(hashedPass)
		} else {
			http.Error(w, "Password must be a string", http.StatusBadRequest)
			return
		}
	}

	updatedUser, err := h.UserRepo.UpdateUser(id, filtered)
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

	json.NewEncoder(w).Encode(updatedUser)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := GetIdFromUrlPath(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.UserRepo.DeleteUser(id)

	if err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User successfully deleted",
	})
}
