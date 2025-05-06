package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/routes"
	"golang.org/x/crypto/bcrypt"
)

func TestUpdateUser(t *testing.T) {
	TearDown(t)

	user := InitUser(t, InitUserOptions{})

	t.Run("Can update user name", func(t *testing.T) {
		updatedName := "updated_first_name"
		body := models.UpdateUserRequest{
			FirstName: &updatedName,
		}

		rr := CreateRouteAndServe(t, routes.UPDATE_USER, ServeOpts{Payload: body, AuthToken: &user.Token, PathParams: map[string]any{"id": user.ID}})

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var updatedUser models.User
		err := json.Unmarshal(rr.Body.Bytes(), &updatedUser)
		if err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		targetUser, err := UserRepo.GetUserById(user.ID)
		if err != nil {
			t.Fatalf("Unable to find user by ID: %v", err)
		}
		if targetUser.FirstName != updatedName {
			t.Fatal("User first name was not updated")
		}
	})

	t.Run("Can update user password", func(t *testing.T) {
		updatedPass := "updated_password"
		body := models.UpdateUserRequest{
			Password:        &updatedPass,
			ConfirmPassword: &updatedPass,
		}

		rr := CreateRouteAndServe(t, routes.UPDATE_USER, ServeOpts{Payload: body, AuthToken: &user.Token, PathParams: map[string]any{"id": user.ID}})

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var updatedUser models.User
		err := json.Unmarshal(rr.Body.Bytes(), &updatedUser)
		if err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		targetUser, err := UserRepo.GetUserById(user.ID)
		if err != nil {
			t.Fatalf("Unable to find user by ID: %v", err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(targetUser.Password), []byte(updatedPass))
		if err != nil {
			t.Fatal("DB password does not match updated password")
			return
		}
	})

	t.Run("Cannot update user password without confirmPassword", func(t *testing.T) {
		updatedPass := "updated_password"
		body := models.UpdateUserRequest{
			Password: &updatedPass,
		}

		rr := CreateRouteAndServe(t, routes.UPDATE_USER, ServeOpts{Payload: body, AuthToken: &user.Token, PathParams: map[string]any{"id": user.ID}})

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("Cannot update user password without matching confirmPassword", func(t *testing.T) {
		updatedPass := "updated_password"
		confPas := "different"
		body := models.UpdateUserRequest{
			Password:        &updatedPass,
			ConfirmPassword: &confPas,
		}

		rr := CreateRouteAndServe(t, routes.UPDATE_USER, ServeOpts{Payload: body, AuthToken: &user.Token, PathParams: map[string]any{"id": user.ID}})

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}
