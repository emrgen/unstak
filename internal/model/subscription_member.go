package model

import "gorm.io/gorm"

type SubscriptionMember struct {
	gorm.Model
	ID             string `gorm:"primaryKey;uuid;not null"`
	SubscriptionID string `gorm:"not null"`
	Subscription   Subscription
	UserID         string `gorm:"not null"`
	Permission     uint64
	CreatedByID    string `gorm:"uuid;not null"`
	UpdateByID     string `gorm:"uuid"`
}
