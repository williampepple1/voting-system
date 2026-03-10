package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type RoleStatus string

const (
	RoleOwner   RoleStatus = "owner"
	RoleAuditor RoleStatus = "auditor"
	RoleVoter   RoleStatus = "voter"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	Username  string     `gorm:"unique" json:"username"`
	Password  string     `json:"-"`
	Zone      string     `json:"zone"`
	Photo     string     `gorm:"type:varchar; null;" json:"photo_url"`
	Role      RoleStatus `gorm:"type:varchar(20);not null;default:'voter'" json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return nil
}
