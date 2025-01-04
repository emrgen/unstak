package model

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&Post{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Course{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Page{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Collection{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&CollectionPost{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Tier{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&TierMember{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Tag{}); err != nil {
		return err
	}

	return nil
}
