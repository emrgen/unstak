package model

type Tag struct {
	ID   string `gorm:"primaryKey;uuid"`
	Name string `gorm:"not null;unique"`
}
