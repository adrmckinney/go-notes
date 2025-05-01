package factories

import (
	"fmt"

	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/utils"
)

// UserFactory generates a slice of models.User for testing or seeding purposes.
// It fetches random user data using utils.GetRandomUsers. If fetching fails,
// the returned users will have zero values for their fields.
//
// Params:
//   - count: the number of users to generate
//   - firstName, lastName, username, password: (ignored) placeholders for compatibility
//
// Returns:
//   - []models.User: a slice of generated User structs with random user data
func UserFactory(count int, firstName string, lastName string, username string, password string) []models.User {
	users := make([]models.User, count)
	rUsers, err := utils.GetRandomUsers(count)
	if err != nil {
		fmt.Printf("failed to fetch random users. Setting to empty string for default value. Error: %v\n", err)
		return users
	}
	for i := 0; i < count; i++ {
		for i := range count {
			users[i] = models.User{
				FirstName: rUsers[i].FirstName,
				LastName:  rUsers[i].LastName,
				Username:  rUsers[i].Username,
				Password:  rUsers[i].Password,
			}
		}

	}
	return users
}
