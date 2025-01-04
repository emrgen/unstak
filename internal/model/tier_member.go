package model

import "gorm.io/gorm"

type TierMember struct {
	gorm.Model
	ID          string `gorm:"primaryKey;uuid;not null"`
	TierID      string `gorm:"not null"`
	Tier        *Tier  `gorm:"foreignKey:TierID;references:ID"`
	UserID      string `gorm:"not null"`
	CreatedByID string `gorm:"uuid;not null"`
	UpdateByID  string `gorm:"uuid"`
}
