package model

import "gorm.io/gorm"

type Reaction struct {
	gorm.Model
	PostID string `gorm:"primaryKey;not null"`
	UserID string `gorm:"primaryKey;not null"`
	Name   string `gorm:"primaryKey;not null"`
	State  bool   `gorm:"not null"`
	Post   *Post  `gorm:"foreignKey:PostID"`
}
