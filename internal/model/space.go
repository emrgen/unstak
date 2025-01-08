package model

import "gorm.io/gorm"

type SpaceConfig struct {
	ConfigDomain      string `json:"config_domain"`
	ConfigCallbackURL string `json:"callback_url"` // redirect the user to this url after login
}

type Space struct {
	gorm.Model
	ID                string      `gorm:"uuid;primaryKey"`
	Name              string      `gorm:"not null;unique"`
	Master            bool        `gorm:"not null;default:false"` // master space is created by the system
	OwnerID           string      `gorm:"uuid;not null"`          // Authbase user id
	Private           bool        `gorm:"not null"`               // data from private spaces are accessed by api key
	PoolID            string      `gorm:"uuid"`                   // for space with separate user pool
	AuthbaseProjectID string      `gorm:"not null"`               // for private spaces, the space admin may choose to keep the user pool separate from the unpost user pool
	SpaceConfig       SpaceConfig `gorm:"embedded"`               // for private spaces, config
}
