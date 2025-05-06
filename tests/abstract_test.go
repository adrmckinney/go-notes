package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/adrmckinney/go-notes/auth"
	"github.com/adrmckinney/go-notes/db"
	"github.com/adrmckinney/go-notes/factories"
	"github.com/adrmckinney/go-notes/handlers"
	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/repos"
	"github.com/adrmckinney/go-notes/routes"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type InitUserOptions struct {
	FirstName string
	Username  string
}

type ServeOpts struct {
	PathParams  map[string]interface{}
	QueryParams map[string]interface{}
	Payload     interface{}
	AuthToken   *string
}

var TestDB *gorm.DB                    // Shared test database connection
var NoteRepo *repos.NoteRepo           // Shared NoteRepo instance
var UserTokenRepo *repos.UserTokenRepo // Shared NoteRepo instance
var UserRepo *repos.UserRepo           // Shared NoteRepo instance
var UserHandler *handlers.UserHandler
var AuthHandler *handlers.AuthHandler
var NoteHandler *handlers.NoteHandler
var router *mux.Router

// var User models.UserWithToken

// TestMain is the entry point for all tests in the tests package and its subpackages
func TestMain(m *testing.M) {
	// Initialize the test database
	TestDB = db.InitTestGorm()

	// Initialize repositories and handlers
	initializeTestDependencies()
	router = routes.NewRouter(TestDB)
	// Run the tests
	code := m.Run()

	// Exit with the test result code
	os.Exit(code)
}

// initializeTestDependencies initializes shared repositories and handlers for tests
func initializeTestDependencies() {
	// Initialize the NoteRepo
	NoteRepo = &repos.NoteRepo{DB: TestDB}
	UserRepo = &repos.UserRepo{DB: TestDB}
	UserTokenRepo = &repos.UserTokenRepo{DB: TestDB}

	// Initialize the NoteHandler with the NoteRepo
	NoteHandler = &handlers.NoteHandler{NoteRepo: NoteRepo}
	UserHandler = &handlers.UserHandler{UserRepo: UserRepo}
	AuthHandler = &handlers.AuthHandler{UserTokenRepo: UserTokenRepo, UserRepo: UserRepo}
}

func TearDown(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		db.CleanUpAllTables(t, TestDB)
	})
}

func InitUser(t *testing.T, opts InitUserOptions) models.UserWithToken {
	t.Helper()
	body := factories.UserFactory(
		factories.UserFactoryOptions{
			Count:     1,
			FirstName: opts.FirstName,
			Username:  opts.Username,
		})

	signupUser := models.SignUpRequest{
		FirstName:       body[0].FirstName,
		LastName:        body[0].LastName,
		Username:        body[0].Username,
		Password:        body[0].Password,
		ConfirmPassword: body[0].Password,
	}

	rr := CreateRouteAndServe(t, routes.SIGN_UP, ServeOpts{Payload: signupUser})

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var user models.UserWithToken
	err := json.Unmarshal(rr.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	// add the password back
	user.Password = signupUser.Password
	return user
}

func InitUsers(t *testing.T, count int, opts InitUserOptions) []models.UserWithToken {
	var users []models.UserWithToken
	for range count {
		user := InitUser(t, opts)
		users = append(users, user)
	}
	return users
}

func CreateRouteAndServe(
	t *testing.T,
	routeDefinition routes.RouteDefinition,
	opts ServeOpts,
) *httptest.ResponseRecorder {
	t.Helper()

	// Replace path parameters in the route
	finalPath := string(routeDefinition.Path)
	for key, value := range opts.PathParams {
		placeholder := fmt.Sprintf("{%s}", key)
		finalPath = strings.ReplaceAll(finalPath, placeholder, fmt.Sprintf("%v", value))
	}

	// Append query parameters
	if len(opts.QueryParams) > 0 {
		values := url.Values{}
		for key, val := range opts.QueryParams {
			values.Set(key, fmt.Sprintf("%v", val))
		}
		finalPath += "?" + values.Encode()
	}

	// Encode request payload if present
	var reqBody io.Reader
	if opts.Payload != nil {
		jsonData, err := json.Marshal(opts.Payload)
		if err != nil {
			t.Fatalf("Failed to encode request payload: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// Create HTTP request
	req, err := http.NewRequest(string(routeDefinition.Method), finalPath, reqBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set headers
	if opts.AuthToken != nil {
		headers := auth.AuthHeaders(*opts.AuthToken)
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	} else if opts.Payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Serve the request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}
