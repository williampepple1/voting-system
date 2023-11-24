package routes

import (
	"voting-system/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func SetupUserRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/register", handlers.RegisterUser(db))
	r.POST("/login", handlers.LoginUser(db))
}
