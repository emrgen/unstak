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
	ID      string `gorm:"primaryKey;uuid"`
	Slug    string `gorm:"not null;unique"`
	Title   string
	Summary string
	Excerpt string
	Content string `gorm:"uuid;not null"`
	//CreatedByID string     `gorm:"not null"`
	Status  PostStatus `gorm:"not null;default:draft"`
	Tags    []*Tag     `gorm:"many2many:post_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Version int64
}

// PostReaction is a map of reaction names to their counts
// this is calculated by aggregating the reactions of all users who reacted to the post from the `reaction` table
type PostReaction map[string]int
