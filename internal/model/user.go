package model

import "gorm.io/gorm"

const (
	UserRoleViewer      = "viewer"
	UserRoleContributor = "contributor"
	UserRoleEditor      = "editor"
	UserRoleAdmin       = "admin"
	UserRoleOwner       = "owner"
)

type User struct {
	gorm.Model
	AuthbaseID string   `gorm:"not null"`
	Role       UserRole `gorm:"not null;default:'viewer'"` // first user to sign up is the owner
}
