package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/adrmckinney/go-notes/routes"
)

func TestLogout(t *testing.T) {
	TearDown(t)

	user := InitUser(t, InitUserOptions{})

	rr := CreateRouteAndServe(t, routes.LOGOUT, ServeOpts{AuthToken: &user.Token})

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resMessage map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resMessage)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	res := UserTokenRepo.TokenExists(user.Token)
	if res {
		t.Fatal("User token was not removed from user_tokens table")
	}
}
