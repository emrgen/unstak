package model

type Tag struct {
	ID      string `gorm:"primaryKey;uuid"`
	SpaceID string `gorm:"not null"`
	Name    string `gorm:"not null"`
}
