package services

import (
	"net/http"

	"github.com/adrmckinney/go-notes/errs"
	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/repos"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepo repos.UserRepo
}

func (s *UserService) GetUserById(userId uint) (models.User, error) {
	user, err := s.UserRepo.GetUserById(userId)
	if err != nil {
		return models.User{}, errs.NewHTTPError(http.StatusNotFound, "User not found")
	}
	return user, nil
}

func (s *UserService) UpdateUser(userId uint, payload map[string]any) (models.User, error) {
	password, passOk := payload["password"].(string)
	confirmPassword, confirmOk := payload["confirmPassword"].(string)

	if passOk && (!confirmOk || confirmPassword == "") {
		return models.User{}, errs.NewHTTPError(http.StatusBadRequest, "ConfirmPassword must be provided for password change")
	}

	if passOk && confirmOk && password != confirmPassword {
		return models.User{}, errs.NewHTTPError(http.StatusBadRequest, "Password and ConfirmPassword do not match")
	}

	filtered := models.FilterUpdateFields(payload, models.AllowedUserUpdateFields)

	if passRaw, ok := filtered["password"]; ok && passRaw != "" {
		if passStr, ok := passRaw.(string); ok {
			hashedPass, err := bcrypt.GenerateFromPassword([]byte(passStr), bcrypt.DefaultCost)
			if err != nil {
				return models.User{}, errs.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
			}
			filtered["password"] = string(hashedPass)
		} else {
			return models.User{}, errs.NewHTTPError(http.StatusBadRequest, "Password must be a string")
		}
	}

	updatedUser, err := s.UserRepo.UpdateUser(userId, filtered)
	if err != nil {
		// Check for GORM's record not found error
		if err == gorm.ErrRecordNotFound ||
			(err.Error() != "" && (err.Error() == "note not found" ||
				err.Error() == "note not found: record not found")) {
			return models.User{}, errs.NewHTTPError(http.StatusNotFound, "Note not found")
		}
		return models.User{}, errs.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return updatedUser, nil
}

func (s *UserService) DeleteUser(userId uint) (map[string]string, error) {
	err := s.UserRepo.DeleteUser(userId)

	if err != nil {
		return map[string]string{}, errs.NewHTTPError(http.StatusInternalServerError, "Failed to delete note")
	}
	return map[string]string{
		"message": "User successfully deleted",
	}, nil
}
