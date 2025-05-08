package main

import (
	"log"
	"net/http"

	"github.com/adrmckinney/go-notes/db"
	"github.com/adrmckinney/go-notes/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("[WARNING] .env file not loaded. Falling back to package environment variables.")
	}

	db.InitGorm()

	r := routes.NewRouter(db.GormDB)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
