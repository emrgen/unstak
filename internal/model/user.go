package model

import "gorm.io/gorm"

type UserRole string

const (
	UserRoleViewer UserRole = "viewer"
	UserRoleEditor          = "editor"
	UserRoleAdmin           = "admin"
	UserRoleOwner           = "owner" // first user who logs-in becomes the owner
)

type User struct {
	gorm.Model
	ID       string   `gorm:"not null"`
	Username string   `gorm:"not null"`
	Email    string   `gorm:"not null"`
	Role     UserRole `gorm:"not null;default:'viewer'"` // first user to sign up is the owner
}
