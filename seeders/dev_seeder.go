package seeders

import (
	"fmt"
	"log"
	"time"

	"github.com/adrmckinney/go-notes/db"
	"github.com/adrmckinney/go-notes/repos"
	"github.com/adrmckinney/go-notes/services"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type SeederDeps struct {
	DB          *gorm.DB
	NoteService *services.NoteService
	UserService *services.UserService
	AuthService *services.AuthService
}

func RunDevSeeders() {
	fmt.Println("Running development seeders...")

	err := godotenv.Load()
	if err != nil {
		log.Println("[WARNING] .env file not loaded. Falling back to package environment variables.")
	}

	fmt.Println("Initializing DB...")
	db.InitGorm()
	fmt.Println("DB initialized")

	// Initialize repo(s)
	userRepo := &repos.UserRepo{DB: db.GormDB}
	userTokenRepo := &repos.UserTokenRepo{DB: db.GormDB}
	noteRepo := &repos.NoteRepo{DB: db.GormDB}

	// Initialize services
	deps := SeederDeps{
		DB:          db.GormDB,
		AuthService: &services.AuthService{UserRepo: *userRepo, UserTokenRepo: *userTokenRepo},
		UserService: &services.UserService{UserRepo: *userRepo},
		NoteService: &services.NoteService{NoteRepo: *noteRepo},
	}

	totalStart := time.Now()

	err = SeedUserNotes(deps, 5, 50)
	if err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	fmt.Printf("All development seeders completed in %v!\n", time.Since(totalStart))
}
