package model

import (
	"gorm.io/gorm"
)

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

type Post struct {
	gorm.Model
	ID           string         `gorm:"primaryKey;uuid"`
	DocumentID   string         `gorm:"uuid;not null"`
	CreatedByID  string         `gorm:"not null"`
	SpaceID      string         `gorm:"uuid;not null"`
	Reaction     string         `gorm:"not null"` // Reaction aggregates the reactions of all users who reacted to the post
	Status       PostStatus     `gorm:"not null;default:draft"`
	Tags         []*Tag         `gorm:"many2many:post_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PlatformTags []*PlatformTag `gorm:"many2many:post_platform_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Authors      []*User        `gorm:"many2many:post_authors;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// PostReaction is a map of reaction names to their counts
// this is calculated by aggregating the reactions of all users who reacted to the post from the `reaction` table
type PostReaction map[string]int
