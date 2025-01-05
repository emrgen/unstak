package model

import (
	"gorm.io/gorm"
	"time"
)

type Page struct {
	gorm.Model
	ID          string  `gorm:"primaryKey;uuid"`
	DocumentID  string  `gorm:"not null"`
	CourseID    string  `gorm:"not null"`
	SpaceID     string  `gorm:"uuid;not null"`
	CreatedByID string  `gorm:"not null"`
	Course      *Course `gorm:"foreignKey:CourseID;references:ID"`
	Status      PostStatus
}

type PageTag struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	PageID    string         `gorm:"not null"`
	TagID     string         `gorm:"not null"`
}
