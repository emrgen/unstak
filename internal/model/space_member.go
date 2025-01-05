package model

import "gorm.io/gorm"

type UserRole string

type SpaceMember struct {
	gorm.Model
	ID      string `gorm:"uuid;primaryKey"`
	UserID  string `gorm:"uuid;not null"`
	SpaceID string `gorm:"uuid;not null"`
	Role    UserRole
}
