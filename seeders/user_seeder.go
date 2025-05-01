package seeders

import (
	"fmt"

	"github.com/adrmckinney/go-notes/db"
	"github.com/adrmckinney/go-notes/factories"
	"github.com/adrmckinney/go-notes/repos"
)

func SeedUsers(count int) {
	userRepo := repos.UserRepo{DB: db.GormDB}

	users := factories.UserFactory(count, "", "", "", "")

	for _, user := range users {
		_, err := userRepo.CreateUser(user)
		if err != nil {
			fmt.Printf("Failed to seed user: %v\n", err)
		}
	}
}
