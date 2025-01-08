package model

import (
	"gorm.io/gorm"
)

type Tier struct {
	gorm.Model
	ID             string `gorm:"primaryKey;uuid"`
	SpaceID        string `gorm:"not null"`
	Name           string `gorm:"not null;uniqueIndex:idx_space_project"`
	CreatedByID    string `gorm:"not null"`
	Free           bool   `gorm:"not null;default:false"`
	MonthlyCost    float64
	YearlyCost     float64
	HalfYearlyCost float64
	QuarterlyCost  float64
}
