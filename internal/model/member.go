package model

import "gorm.io/gorm"

type OutletMember struct {
	gorm.Model
	OutletID   string `gorm:"primaryKey;not null"`
	UserID     string `gorm:"primaryKey;not null"`
	Permission uint64
	OwnerID    string `gorm:"uuid;not null"`
	UpdateByID string `gorm:"uuid"`
}
