package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/adrmckinney/go-notes/db"
	"github.com/adrmckinney/go-notes/routes"
	"github.com/adrmckinney/go-notes/seeders"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Define the flags
	seedDev := flag.Bool("seedDev", false, "Run development seeders")
	migrateFlag := flag.Bool("migrate", false, "Run database migrations")
	rollbackFlag := flag.Bool("migrate:rollback", false, "Rollback database migrations")
	stepsFlag := flag.Int("steps", 0, "Number of migration steps to rollback (default: 0 for all)")
	removeFlag := flag.Bool("remove", false, "Removes the entire mckinney_go_notes_db")
	flag.Parse()

	db.Init()

	if *seedDev {
		seeders.RunDevSeeders()
		return
	}
	if *migrateFlag {
		db.RunMigrations()
		return
	}
	if *rollbackFlag {
		db.RunRollback(*stepsFlag)
		return
	}
	if *removeFlag {
		db.RunRemoveDatabase()
		return
	}

	r := routes.NewRouter(db.DB)
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
