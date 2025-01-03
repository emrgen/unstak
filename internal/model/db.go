package model

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&Post{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Subscription{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&SubscriptionMember{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Tag{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Course{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Page{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&PageTag{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Collection{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&CollectionPost{}); err != nil {
		return err
	}

	return nil
}
