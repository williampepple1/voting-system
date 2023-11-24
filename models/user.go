package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type RoleStatus string

const (
	StatusAdmin   RoleStatus = "owner"
	StatusAuditor RoleStatus = "auditor"
	StatusVoter   RoleStatus = "voter"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Username  string     `gorm:"unique" json:"username"`
	Password  string     `json:"password"`
	Zone      string     `json:"zone"`
	Photo     string     `gorm:"type:varchar; null;" json:"photo_url"`
	Role      RoleStatus `gorm:"type:varchar(20);not null;default:'voter'" json:"role"`
	// Role
	// Tasks     []Task    `gorm:"foreignkey:UserID"` // This indicates a one-to-many relationship
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *User) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID = uuid.New()
	currentTime := time.Now()
	base.CreatedAt = currentTime
	base.UpdatedAt = currentTime
	return nil
}

// GORM V2 uses callbacks like BeforeUpdate to handle the update timestamp
func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	user.UpdatedAt = time.Now()
	return
}
