package seeders

import (
	"fmt"
	"time"
)

func RunDevSeeders() {
	fmt.Println("Running development seeders...")

	totalStart := time.Now()

	_ = SeedUserNotes(5, 50)

	fmt.Printf("All development seeders completed in %v!\n", time.Since(totalStart))
}
