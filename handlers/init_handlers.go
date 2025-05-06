package handlers

import (
	"github.com/adrmckinney/go-notes/repos"
	"gorm.io/gorm"
)

type Handlers struct {
	AuthHandler *AuthHandler
	UserHandler *UserHandler
	NoteHandler *NoteHandler
}

func InitHandlers(db *gorm.DB) *Handlers {
	userRepo := &repos.UserRepo{DB: db}
	userTokenRepo := &repos.UserTokenRepo{DB: db}
	noteRepo := &repos.NoteRepo{DB: db}

	// Initialize handlers with repositories
	authHandler := &AuthHandler{
		UserRepo:      userRepo,
		UserTokenRepo: userTokenRepo,
	}
	userHandler := &UserHandler{
		UserRepo: userRepo,
	}
	noteHandler := &NoteHandler{
		NoteRepo: noteRepo,
	}

	return &Handlers{
		AuthHandler: authHandler,
		UserHandler: userHandler,
		NoteHandler: noteHandler,
	}
}
