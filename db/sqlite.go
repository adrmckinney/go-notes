package db

import (
	"log"

	"github.com/adrmckinney/go-notes/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func InitTestGorm() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to sqlite: %v", err)
	}
	db.AutoMigrate(&models.Note{})
	return db
}
