package model

import "gorm.io/gorm"

// Collection represents a collection of posts by a user
type Collection struct {
	gorm.Model
	ID          string `gorm:"primaryKey;uuid"`
	Name        string `gorm:"not null"`
	CreatedByID string `gorm:"uuid;not null"`
	Private     bool   `gorm:"default:true"`
	Tags        []*Tag `gorm:"many2many:collection_tags;"` // Tags are automatically added to the collection when a post is added
}

type CollectionTag struct {
	gorm.Model
	CollectionID string `gorm:"primaryKey;not null"`
	TagID        string `gorm:"primaryKey;not null"`
}
