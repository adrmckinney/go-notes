package main

import (
	"log"
	"net/http"

	"github.com/adrmckinney/go-notes/db"
	"github.com/adrmckinney/go-notes/routes"
)

func main() {
	db.InitGorm()

	r := routes.NewRouter(db.GormDB)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
