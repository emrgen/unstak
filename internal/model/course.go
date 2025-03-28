package model

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	ID           string         `gorm:"primaryKey;uuid"`
	DocumentID   string         `gorm:"not null"`
	CreatedByID  string         `gorm:"not null"`
	SpaceID      string         `gorm:"uuid;not null"`
	Status       PostStatus     `gorm:"not null;default:draft"`
	Reaction     string         `gorm:"not null"` // Reaction aggregates the reactions of all users who reacted to the post
	Tags         []*Tag         `gorm:"many2many:course_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PlatformTags []*PlatformTag `gorm:"many2many:course_platform_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Authors      []*User        `gorm:"many2many:course_authors;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
