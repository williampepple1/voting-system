package main

import (
	"log"
	"os"
	"voting-system/config"
	"voting-system/models"
	"voting-system/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading environment variables from system")
	}

	cfg, err := config.LoadAppConfig()
	if err != nil {
		log.Printf("Configuration error: %v", err)
		os.Exit(1)
	}

	db, err := config.InitDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.AutoMigrate(&models.User{}).Error; err != nil {
		log.Printf("Failed to migrate database schema: %v", err)
		os.Exit(1)
	}

	r := gin.Default()
	routes.SetupUserRoutes(r, db, cfg)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Printf("Server failed: %v", err)
		os.Exit(1)
	}
}
