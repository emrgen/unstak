package model

import "gorm.io/gorm"

type SpaceRole string

const (
	SpaceRoleViewer      SpaceRole = "viewer"
	SpaceRoleContributor           = "contributor"
	SpaceRoleEditor                = "editor"
	SpaceRoleAdmin                 = "admin"
	SpaceRoleOwner                 = "owner"
)

type SpaceMember struct {
	gorm.Model
	UserID  string `gorm:"uuid;primaryKey;not null;"`
	SpaceID string `gorm:"uuid;primaryKey;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role    UserRole
	Space   *Space
	User    *User
}
