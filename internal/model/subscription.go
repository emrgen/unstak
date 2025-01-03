package model

import (
	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	ID          string `gorm:"primaryKey;uuid"`
	Name        string `gorm:"not null;uniqueIndex:idx_space_project"`
	CreatedByID string `gorm:"not null"`
	UserDefault bool
}
