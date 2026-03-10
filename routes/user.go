package routes

import (
	"voting-system/config"
	"voting-system/handlers"
	"voting-system/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func SetupUserRoutes(r *gin.Engine, db *gorm.DB, cfg *config.AppConfig) {
	// Public routes
	r.POST("/register", handlers.RegisterUser(db))
	r.POST("/login", handlers.LoginUser(db, cfg))

	// Protected routes — require a valid JWT
	protected := r.Group("/")
	protected.Use(middleware.Authorize(cfg))
	{
		protected.GET("/me", handlers.GetCurrentUser(db))
	}
}
