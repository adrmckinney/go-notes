package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() {
	var err error
	dsn := "root:password@tcp(127.0.0.1:3306)/go_notes_dev"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database unreachable:", err)
	}
}
