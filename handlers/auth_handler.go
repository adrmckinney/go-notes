package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/repos"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo      *repos.UserRepo
	UserTokenRepo *repos.UserTokenRepo
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

	existingUser, err := h.UserRepo.GetUserByUsername(req.Username)
	if err == nil && existingUser.ID != 0 {
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Password:  string(hashedPass),
	}

	createdUser, err := h.UserRepo.CreateUser(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Do not return the password in the response
	createdUser.Password = ""
	// Create token and log user in
	tokenString, expiry, err := generateAuthToken(createdUser)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	userToken := models.UserToken{
		UserId:    createdUser.ID,
		Token:     tokenString,
		ExpiresAt: expiry,
	}

	err = h.UserTokenRepo.StoreUserToken(userToken)
	if err != nil {
		http.Error(w, "Failed to store token", http.StatusInternalServerError)
		return
	}

	// Add token to response
	userWithToken := models.UserWithToken{
		User:  createdUser,
		Token: tokenString,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userWithToken)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req models.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	targetUser, err := h.UserRepo.GetUserByUsername(*req.Username)
	if err != nil {
		http.Error(w, "Username or password are incorrect", http.StatusForbidden)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(targetUser.Password), []byte(*req.Password))
	if err != nil {
		http.Error(w, "Username or password are incorrect", http.StatusForbidden)
		return
	}

	tokenString, expiry, err := generateAuthToken(targetUser)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	userToken := models.UserToken{
		UserId:    targetUser.ID,
		Token:     tokenString,
		ExpiresAt: expiry,
	}

	err = h.UserTokenRepo.StoreUserToken(userToken)
	if err != nil {
		http.Error(w, "Failed to store token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Shouldn't even get here because of AuthMiddleware
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	var tokenString string
	fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
	if tokenString == "" {
		// Shouldn't even get here because of AuthMiddleware
		http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
		return
	}

	// Remove token from user_tokens table
	err := h.UserTokenRepo.DeleteUserToken(tokenString)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func generateAuthToken(user models.User) (string, time.Time, error) {
	var jwtKey = []byte("")

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate token: %w", err)

	}
	return tokenString, expirationTime, nil
}
