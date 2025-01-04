package model

import "gorm.io/gorm"

type NewsLetter struct {
	gorm.Model
	ID     string `gorm:"primaryKey;uuid"`
	UserID string `gorm:"not null"`
	Active bool   `gorm:"not null"`
	Emails string `gorm:"not null"` // comma separated emails
}
