package seeders

import (
	"github.com/adrmckinney/go-notes/factories"
	"github.com/adrmckinney/go-notes/models"
)

func SeedUsers(deps SeederDeps, count int) ([]uint, error) {
	users := factories.UserFactory(factories.UserFactoryOptions{Count: count})
	var userIds []uint
	for i, user := range users {
		if i == 0 {
			user.Username = "demoUser"
			user.Password = "password"
		}
		signupUser := models.SignUpRequest{
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Username:        user.Username,
			Password:        user.Password,
			ConfirmPassword: user.Password,
		}
		created, err := deps.AuthService.SignUp(signupUser)
		if err != nil {
			return nil, err
		}
		userIds = append(userIds, created.ID)
	}

	return userIds, nil
}
