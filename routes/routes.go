package routes

import (
	"os"

	"github.com/adrmckinney/go-notes/handlers"
	"github.com/adrmckinney/go-notes/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type HTTPMethod string
type Route string

type RouteDefinition struct {
	Method HTTPMethod
	Path   Route
}

const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
)

var (
	SIGN_UP     = RouteDefinition{Method: POST, Path: "/signup"}
	SIGN_IN     = RouteDefinition{Method: POST, Path: "/signin"}
	LOGOUT      = RouteDefinition{Method: DELETE, Path: "/logout"}
	UPDATE_USER = RouteDefinition{Method: PUT, Path: "/users/{id}"}
	GET_NOTES   = RouteDefinition{Method: GET, Path: "/notes"}
	GET_NOTE    = RouteDefinition{Method: GET, Path: "/notes/{id}"}
	CREATE_NOTE = RouteDefinition{Method: POST, Path: "/notes"}
	UPDATE_NOTE = RouteDefinition{Method: PUT, Path: "/notes/{id}"}
	DELETE_NOTE = RouteDefinition{Method: DELETE, Path: "/notes/{id}"}
)

func NewRouter(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()
	var jwtKey = []byte(os.Getenv("JWT_SECRET")) // TODO this needs to be added to DOCKER

	// Implemented init handlers because testing requires a
	// sqlite db to be created and used, which means we need to
	// pass in the correct DB here.
	handlers := handlers.InitHandlers(db)

	// Guest routes
	r.HandleFunc(string(SIGN_UP.Path), handlers.AuthHandler.SignUp).Methods(string(SIGN_UP.Method))
	r.HandleFunc(string(SIGN_IN.Path), handlers.AuthHandler.SignIn).Methods(string(SIGN_UP.Method))

	// Auth routes
	auth := r.PathPrefix("/").Subrouter()
	auth.Use(middleware.AuthMiddleware(jwtKey, *handlers.AuthHandler.UserTokenRepo))

	// User
	auth.HandleFunc(string(LOGOUT.Path), handlers.AuthHandler.Logout).Methods(string(LOGOUT.Method))
	auth.HandleFunc(string(UPDATE_USER.Path), handlers.UserHandler.UpdateUser).Methods(string(UPDATE_USER.Method))

	// Notes
	auth.HandleFunc(string(GET_NOTES.Path), handlers.NoteHandler.GetNotes).Methods(string(GET_NOTES.Method))
	auth.HandleFunc(string(GET_NOTE.Path), handlers.NoteHandler.GetNote).Methods(string(GET_NOTE.Method))
	auth.HandleFunc(string(CREATE_NOTE.Path), handlers.NoteHandler.CreateNote).Methods(string(CREATE_NOTE.Method))
	auth.HandleFunc(string(UPDATE_NOTE.Path), handlers.NoteHandler.UpdateNote).Methods(string(UPDATE_NOTE.Method))
	auth.HandleFunc(string(DELETE_NOTE.Path), handlers.NoteHandler.DeleteNote).Methods(string(DELETE_NOTE.Method))
	return r
}
