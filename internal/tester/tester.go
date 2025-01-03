package tester

import (
	"os"

	"github.com/emrgen/unpost/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	testPath = "../../.test/"
)

var (
	db *gorm.DB
)

func Setup() {
	RemoveDBFile()

	_ = os.Setenv("ENV", "test")

	err := os.MkdirAll(testPath+"/db", os.ModePerm)
	if err != nil {
		panic(err)
	}

	db, err = gorm.Open(sqlite.Open(testPath+"db/unpost.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = model.Migrate(db)
	if err != nil {
		panic(err)
	}
}

func TestDB() *gorm.DB {
	return db
}

func RemoveDBFile() {
	err := os.RemoveAll(testPath)
	if err != nil {
		panic(err)
	}
}

func CleanUp() {
	DropDocuments()
	DropGroups()
	DropProjects()
	DropUsers()
}

func DropDocuments() {
	db.Exec("DELETE FROM documents")
}

func DropUsers() {
	db.Exec("DELETE FROM users")
}

func DropProjects() {
	db.Exec("DELETE FROM projects")
}

func DropGroups() {
	db.Exec("DELETE FROM groups")
}
