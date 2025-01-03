package model

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	ID          string `gorm:"primaryKey;not null"`
	Name        string `gorm:"uniqueIndex;not null"`
	CreatedByID string `gorm:"uuid;not null"`
	UpdateByID  string `gorm:"uuid"`
}
