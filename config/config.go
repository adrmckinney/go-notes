package config

import (
	"os"
)

// DBConfig holds the database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func GetDBConfig() DBConfig {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}

	database := os.Getenv("DB_NAME")
	if database == "" {
		database = "mckinney_go_notes_db"
	}

	return DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
	}
}

// GetDSN returns the Data Source Name (DSN) for MySQL (includes db name)
func GetDSN() string {
	config := GetDBConfig()
	return config.User + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.Database + "?multiStatements=true"
}

// GetServerDSN returns the DSN for connecting to the MySQL server (without specifying a database)
func GetServerDSN() string {
	config := GetDBConfig()
	return config.User + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/"
}
