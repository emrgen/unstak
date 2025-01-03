package model

import "gorm.io/gorm"

type CollectionPost struct {
	gorm.Model
	CollectionID string `gorm:"not null"`
	PostID       string `gorm:"not null"`
}
