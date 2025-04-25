package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/adrmckinney/go-notes/db"
	"github.com/adrmckinney/go-notes/routes"
	"github.com/adrmckinney/go-notes/seeders"
)

func main() {
	// Define the seedDev flag
	seedDev := flag.Bool("seedDev", false, "Run development seeders")
	flag.Parse()

	db.Init()

	if *seedDev {
		seeders.RunDevSeeders()
		return
	}

	r := routes.NewRouter(db.DB)
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
