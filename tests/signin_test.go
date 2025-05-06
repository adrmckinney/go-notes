package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/routes"
)

func TestSignin(t *testing.T) {
	TearDown(t)

	user := InitUser(t, InitUserOptions{})

	body := models.SignInRequest{
		Username: &user.Username,
		Password: &user.Password,
	}

	rr := CreateRouteAndServe(t, routes.SIGN_IN, ServeOpts{Payload: body})

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var res map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	if res["token"] == "" {
		t.Fatalf("Token was not returned in response: %v", err)
	}
}
