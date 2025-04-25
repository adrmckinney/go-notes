package seeders

import (
	"fmt"
	"time"
)

func RunDevSeeders() {
	fmt.Println("Running development seeders...")

	// Call each seeder with a predefined count
	start := time.Now()
	fmt.Println("Seeding notes...")
	SeedNotes(10)
	fmt.Printf("Notes seed in %v\n", time.Since(start))

	// Add other seeders here with timing
	// Example:
	// start = time.Now()
	// fmt.Println("Seeding users...")
	// SeedUsers(5) // Seed 5 users
	// fmt.Printf("Users seeded in %v\n", time.Since(start))

	fmt.Println("Development seeders completed!")
}
