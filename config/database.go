package config

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// InitDB initializes and returns a connection to the PostgreSQL database.
// Environment variables must be loaded before calling this function.
func InitDB() (*gorm.DB, error) {
	host := os.Getenv("DATABASE_HOST")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	port := os.Getenv("DATABASE_PORT")
	if port == "" {
		port = "5432"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		host, port, user, dbname, password,
	)

	return gorm.Open("postgres", connStr)
}
