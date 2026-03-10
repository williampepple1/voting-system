package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
	"voting-system/config"
	"voting-system/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const tokenExpiry = 15 * time.Minute

func RegisterUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		lowerUsername := strings.ToLower(user.Username)

		var existing models.User
		result := db.Where("lower(username) = ?", lowerUsername).First(&existing)
		if result.Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
			return
		}
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("RegisterUser: DB error checking username: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing username"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("RegisterUser: bcrypt error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
		user.Password = string(hashedPassword)
		user.Username = lowerUsername

		if err := db.Create(&user).Error; err != nil {
			log.Printf("RegisterUser: DB create error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "id": user.ID})
	}
}

func LoginUser(db *gorm.DB, cfg *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.User
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		lowerUsername := strings.ToLower(input.Username)

		var user models.User
		if err := db.Where("username = ?", lowerUsername).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Username or password incorrect"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Username or password incorrect"})
			return
		}

		now := time.Now()
		claims := jwt.RegisteredClaims{
			Subject:   user.Username,
			ID:        user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenExpiry)),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(cfg.JWTSecretKey))
		if err != nil {
			log.Printf("LoginUser: token signing error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

func GetCurrentUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var user models.User
		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			log.Printf("GetCurrentUser: DB error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"zone":       user.Zone,
			"photo_url":  user.Photo,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		})
	}
}
