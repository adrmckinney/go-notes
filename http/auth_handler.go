package http

import (
	"encoding/json"
	"net/http"

	"github.com/adrmckinney/go-notes/auth"
	"github.com/adrmckinney/go-notes/errs"
	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/services"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	AuthService services.AuthService
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req models.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(req); err != nil {
		http.Error(w, "All fields must be completed", http.StatusBadRequest)
		return
	}

	result, err := h.AuthService.SignUp(req)
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

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var payload models.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(payload); err != nil {
		http.Error(w, "All fields must be completed", http.StatusBadRequest)
		return
	}

	result, err := h.AuthService.SignIn(payload)
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

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	authInfo, ok := auth.GetAuthInfo(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Remove token from user_tokens table
	res, err := h.AuthService.Logout(authInfo.Token)
	if err != nil {
		if httpErr, ok := err.(*errs.HTTPError); ok {
			http.Error(w, httpErr.Message, httpErr.Status)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
