package model

import "gorm.io/gorm"

// PlatformTag is a model for platform tags
// A platform tag is a tag that is used to categorize and create suggestions in the platform
type PlatformTag struct {
	gorm.Model
	ID   string `gorm:"primaryKey;uuid"`
	Name string `gorm:"not null;unique"`
}
