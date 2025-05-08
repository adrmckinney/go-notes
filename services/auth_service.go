package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/adrmckinney/go-notes/config"
	"github.com/adrmckinney/go-notes/errs"
	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/repos"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo      repos.UserRepo
	UserTokenRepo repos.UserTokenRepo
}

func (s *AuthService) SignIn(payload models.SignInRequest) (map[string]string, error) {
	targetUser, err := s.UserRepo.GetUserByUsername(*payload.Username)
	if err != nil {
		return map[string]string{}, errs.NewHTTPError(http.StatusForbidden, "Username or password are incorrect")
	}

	err = bcrypt.CompareHashAndPassword([]byte(targetUser.Password), []byte(*payload.Password))
	if err != nil {
		return map[string]string{}, errs.NewHTTPError(http.StatusForbidden, "Username or password are incorrect")
	}

	tokenString, expiry, err := generateAuthToken(targetUser)
	if err != nil {
		return map[string]string{}, errs.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	userToken := models.UserToken{
		UserId:    targetUser.ID,
		Token:     tokenString,
		ExpiresAt: expiry,
	}

	err = s.UserTokenRepo.StoreUserToken(userToken)
	if err != nil {
		return map[string]string{}, errs.NewHTTPError(http.StatusInternalServerError, "Failed to store token")
	}

	return map[string]string{"token": tokenString}, nil
}

func (s *AuthService) SignUp(req models.SignUpRequest) (models.UserWithToken, error) {
	// Validate
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return models.UserWithToken{}, fmt.Errorf("validation error: %w", err)
	}

	// Check for existing
	existingUser, err := s.UserRepo.GetUserByUsername(req.Username)
	if err == nil && existingUser.ID != 0 {
		return models.UserWithToken{}, errs.NewHTTPError(http.StatusConflict, "Username already exists")
	}

	// Hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserWithToken{}, errs.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to hash password: %w", err).Error())
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Password:  string(hashedPass),
	}

	createdUser, err := s.UserRepo.CreateUser(user)
	if err != nil {
		return models.UserWithToken{}, fmt.Errorf("failed to create user: %w", err)
	}

	createdUser.Password = ""

	tokenString, expiry, err := generateAuthToken(createdUser)
	if err != nil {
		return models.UserWithToken{}, fmt.Errorf("failed to generate token: %w", err)
	}

	token := models.UserToken{
		UserId:    createdUser.ID,
		Token:     tokenString,
		ExpiresAt: expiry,
	}

	if err := s.UserTokenRepo.StoreUserToken(token); err != nil {
		return models.UserWithToken{}, fmt.Errorf("failed to store token: %w", err)
	}

	return models.UserWithToken{
		User:  createdUser,
		Token: tokenString,
	}, nil
}

func (s *AuthService) Logout(authToken string) (map[string]string, error) {
	err := s.UserTokenRepo.DeleteUserToken(authToken)
	if err != nil {
		return map[string]string{}, errs.NewHTTPError(http.StatusInternalServerError, "Failed to logout")
	}
	return map[string]string{"message": "Logged out successfully"}, nil
}

func generateAuthToken(user models.User) (string, time.Time, error) {
	var jwtKey = []byte(config.GetConfig().JwtSecret)

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
