package model

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	ID          string `gorm:"primaryKey;uuid"`
	Title       string `gorm:"not null"`
	ContentPage string `gorm:"not null"`
	CreatedByID string `gorm:"not null"`
}

type BookTag struct {
	gorm.Model
	BookID string `gorm:"not null"`
	TagID  string `gorm:"not null"`
}
