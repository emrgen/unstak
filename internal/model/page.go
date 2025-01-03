package model

import (
	"gorm.io/gorm"
	"time"
)

type Page struct {
	gorm.Model
	ID           string `gorm:"primaryKey;uuid"`
	Title        string `gorm:"not null"`
	ContentPages string `gorm:"not null"`
	BookID       string `gorm:"not null"`
	Book         *Book  `gorm:"foreignKey:BookID;references:ID"`
}

type PageTag struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	PageID    string         `gorm:"not null"`
	TagID     string         `gorm:"not null"`
}
