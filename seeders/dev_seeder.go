package seeders

import (
	"fmt"
	"time"
)

func RunDevSeeders() {
	fmt.Println("Running development seeders...")

	totalStart := time.Now()

	start := time.Now()
	fmt.Println("Seeding users...")
	SeedUsers(2)
	fmt.Printf("Users seeded in %v\n", time.Since(start))

	start = time.Now()
	fmt.Println("Seeding notes...")
	SeedNotes(10)
	fmt.Printf("Notes seeded in %v\n", time.Since(start))

	fmt.Printf("All development seeders completed in %v!\n", time.Since(totalStart))
}
