package model

import "gorm.io/gorm"

type UserRole string

const (
	UserRoleViewer      = "viewer"
	UserRoleContributor = "contributor"
	UserRoleEditor      = "editor"
	UserRoleAdmin       = "admin"
	UserRoleOwner       = "owner"
)

type User struct {
	gorm.Model
	ID         string   `gorm:"uuid;primaryKey;"`
	AuthbaseID string   `gorm:"not null"`
	Role       UserRole `gorm:"not null"`
}
