package model

import "gorm.io/gorm"

type Space struct {
	gorm.Model
	ID          string `gorm:"uuid;primaryKey"`
	CreatedByID string `gorm:"uuid;not null"`
	Private     bool   `gorm:"not null"` // data from private spaces are accessed by api key
}
