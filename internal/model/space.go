package model

import (
	"gorm.io/gorm"
)

type Outlet struct {
	gorm.Model
	ID          string `gorm:"primaryKey;uuid"`
	ProjectID   string `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;uniqueIndex:idx_space_project"`
	Name        string `gorm:"not null;uniqueIndex:idx_space_project"`
	OwnerID     string `gorm:"not null"`
	UserDefault bool
}
