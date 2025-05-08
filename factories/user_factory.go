package factories

import (
	"fmt"

	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/utils"
)

type UserFactoryOptions struct {
	Count     int
	FirstName string
	LastName  string
	Username  string
	Password  string
}

// UserFactory generates a slice of models.User for testing or seeding purposes.
// It fetches random user data using utils.GetRandomUsers. If fetching fails,
// the returned users will have zero values for their fields.
//
// Params:
//   - count: the number of users to generate. If Count is not passed in, the func will default to 1
//   - firstName, lastName, username, password: (ignored) placeholders for compatibility
//
// Returns:
//   - []models.User: a slice of generated User structs with random user data
//
// func UserFactory(count int, firstName string, lastName string, username string, password string) []models.User {
func UserFactory(opts UserFactoryOptions) []models.User {
	if opts.Username != "" && opts.Count > 1 {
		panic("UserFactory: cannot generate multiple users with the same username")
	}

	if opts.Count == 0 {
		opts.Count = 1
	}

	users := make([]models.User, opts.Count)
	// Get random users only if needed
	var rUsers []utils.RandomUser
	var err error
	if opts.FirstName == "" || opts.LastName == "" || opts.Username == "" || opts.Password == "" {
		rUsers, err = utils.GetRandomUsers(opts.Count)
		if err != nil {
			fmt.Printf("failed to fetch random users. Setting default empty values. Error: %v\n", err)
		}
	}

	for i := range opts.Count {
		user := models.User{
			FirstName: opts.FirstName,
			LastName:  opts.LastName,
			Username:  opts.Username,
			Password:  opts.Password,
		}

		// Only fill with random values if field was not passed in
		if user.FirstName == "" && len(rUsers) > i {
			user.FirstName = rUsers[i].FirstName
		}
		if user.LastName == "" && len(rUsers) > i {
			user.LastName = rUsers[i].LastName
		}
		if user.Username == "" && len(rUsers) > i {
			user.Username = rUsers[i].Username
		}
		if user.Password == "" && len(rUsers) > i {
			user.Password = rUsers[i].Password
		}

		users[i] = user
	}
	return users
}
