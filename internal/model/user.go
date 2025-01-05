package model

import "gorm.io/gorm"

type UserRole string

const (
	UserRoleViewer      UserRole = "viewer"
	UserRoleContributor          = "contributor"
	UserRoleEditor               = "editor"
	UserRoleAdmin                = "admin"
	UserRoleOwner                = "owner"
)

type User struct {
	gorm.Model
	ID      string   `gorm:"not null"`
	SpaceID string   // used for space user pool
	Role    UserRole `gorm:"not null;default:'viewer'"` // first user to sign up is the owner
}
