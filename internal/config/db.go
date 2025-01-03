package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"time"
)

func GetDb(config *Config) *gorm.DB {
	var rdb *gorm.DB
	var err error

	if config.DbConfig.Type == "postgres" {
		rdb, err = gorm.Open(postgres.Open(config.DbConfig.ConnectionString), &gorm.Config{})
	} else {
		filePath := os.Getenv("SQLITE_FILE_PATH")
		if filePath == "" {
			filePath = ".tmp/db/_document.db"
		}

		rdb, err = gorm.Open(sqlite.Open(filePath), &gorm.Config{})
	}

	if err != nil {
		panic(err)
	}

	if config.DbConfig.Type == "postgres" {
		pdb, err := rdb.DB()
		if err != nil {
			panic(err)
		}

		//pdb.SetMaxOpenConns(100)
		//pdb.SetMaxIdleConns(10)
		pdb.SetConnMaxLifetime(time.Hour)
	}

	return rdb
}
