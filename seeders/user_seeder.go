package seeders

import (
	"github.com/adrmckinney/go-notes/db"
	"github.com/adrmckinney/go-notes/factories"
	"github.com/adrmckinney/go-notes/repos"
)

func SeedUsers(count int) ([]uint, error) {
	userRepo := repos.UserRepo{DB: db.GormDB}

	users := factories.UserFactory(factories.UserFactoryOptions{Count: count})
	var userIds []uint
	for _, user := range users {
		created, err := userRepo.CreateUser(user)
		if err != nil {
			return nil, err
		}
		userIds = append(userIds, created.ID)
	}

	return userIds, nil
}
