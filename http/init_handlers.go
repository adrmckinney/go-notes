package http

import (
	"github.com/adrmckinney/go-notes/repos"
	"github.com/adrmckinney/go-notes/services"
	"gorm.io/gorm"
)

type Handlers struct {
	AuthHandler *AuthHandler
	UserHandler *UserHandler
	NoteHandler *NoteHandler
}

func InitHandlers(db *gorm.DB) *Handlers {
	// Initialize Repositories
	userRepo := &repos.UserRepo{DB: db}
	userTokenRepo := &repos.UserTokenRepo{DB: db}
	noteRepo := &repos.NoteRepo{DB: db}

	// Initialize Services with Repos
	authService := &services.AuthService{
		UserRepo:      *userRepo,
		UserTokenRepo: *userTokenRepo,
	}
	userService := &services.UserService{
		UserRepo: *userRepo,
	}
	noteService := &services.NoteService{
		NoteRepo: *noteRepo,
	}

	// Initialize Handlers with Services
	authHandler := &AuthHandler{
		AuthService: *authService,
	}
	userHandler := &UserHandler{
		UserService: *userService,
	}
	noteHandler := &NoteHandler{
		NoteService: *noteService,
	}

	return &Handlers{
		AuthHandler: authHandler,
		UserHandler: userHandler,
		NoteHandler: noteHandler,
	}
}
