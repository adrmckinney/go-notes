package main

import (
	"github.com/adrmckinney/go-notes/db"
	"github.com/adrmckinney/go-notes/seeders"
)

func main() {
	db.InitGorm()
	seeders.RunDevSeeders()
}
