package model

import "gorm.io/gorm"

type Organization struct {
	gorm.Model
	ID   string `gorm:"primaryKey;uuid"`
	Name string `gorm:"not null"`
}
