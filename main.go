package main

import (
	"voting-system/config"
	"voting-system/models"
	"voting-system/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db, err := config.InitDB() // Initialize the database connection
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()

	// Migrate the schema

	db.AutoMigrate(&models.User{})

	// Set up routes
	routes.SetupUserRoutes(r, db) // Setup user routes

	r.Run(":8080")
}
