package model

type Tag struct {
	ID      string `gorm:"primaryKey;uuid"`
	SpaceID string `gorm:"not null;uniqueIndex:idx_space_name"`
	Name    string `gorm:"not null;uniqueIndex:idx_space_name"`
}
