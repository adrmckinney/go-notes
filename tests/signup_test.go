package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/routes"
)

func TestSignup(t *testing.T) {
	TearDown(t)

	body := models.SignUpRequest{
		FirstName:       "testFirst",
		LastName:        "testLast",
		Username:        "testUser",
		Password:        "password",
		ConfirmPassword: "password",
	}

	rr := CreateRouteAndServe(t, routes.SIGN_UP, ServeOpts{Payload: body})

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var createdUser models.UserWithToken
	err := json.Unmarshal(rr.Body.Bytes(), &createdUser)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Verify the token and user were inserted into their databases
	tokenExists := UserTokenRepo.TokenExists(createdUser.Token)
	if !tokenExists {
		t.Fatalf("Token does not exist in user_tokens table: %v", err)
	}

	_, err = UserRepo.GetUserByUsername(createdUser.Username)
	if err != nil {
		t.Fatalf("Username does not exist in users table: %v", err)
	}
}
