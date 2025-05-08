package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	stdhttp "net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/adrmckinney/go-notes/auth"
	"github.com/adrmckinney/go-notes/db"
	"github.com/adrmckinney/go-notes/factories"
	myhttp "github.com/adrmckinney/go-notes/http"
	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/repos"
	"github.com/adrmckinney/go-notes/routes"
	"github.com/adrmckinney/go-notes/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type InitUserOptions struct {
	FirstName string
	Username  string
}

type ServeOpts struct {
	PathParams  map[string]any
	QueryParams map[string]any
	Payload     any
	AuthToken   *string
}

var TestDB *gorm.DB // Shared test database connection
var NoteRepo *repos.NoteRepo
var UserTokenRepo *repos.UserTokenRepo
var UserRepo *repos.UserRepo

// Shared services
var NoteService *services.NoteService
var UserService *services.UserService
var AuthService *services.AuthService

// Shared HTTP handlers
var NoteHandler *myhttp.NoteHandler
var UserHandler *myhttp.UserHandler
var AuthHandler *myhttp.AuthHandler

var router *mux.Router

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("[WARNING] .env file not loaded. Falling back to package environment variables.")
	}

	// Initialize the test database
	TestDB = db.InitTestGorm()

	// Initialize repositories, services, and handlers and set test variables
	handlers := myhttp.InitHandlers(TestDB)
	AuthHandler = handlers.AuthHandler
	UserHandler = handlers.UserHandler
	NoteHandler = handlers.NoteHandler

	AuthService = &AuthHandler.AuthService
	UserService = &UserHandler.UserService
	NoteService = &NoteHandler.NoteService

	UserTokenRepo = &AuthService.UserTokenRepo
	UserRepo = &UserService.UserRepo
	NoteRepo = &NoteService.NoteRepo

	// Initialize router with TestDB
	router = routes.NewRouter(TestDB)

	// Run the tests
	code := m.Run()

	// Exit with the test result code
	os.Exit(code)
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

	user, err := AuthService.SignUp(signupUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err.Error())
	}

	// add the password back because service removes it from response
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
	req, err := stdhttp.NewRequest(string(routeDefinition.Method), finalPath, reqBody)
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
