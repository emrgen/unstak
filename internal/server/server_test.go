package server

import (
	"github.com/emrgen/unpost/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"testing"
)

const (
	testPath = "../../.test/"
)

var (
	db *gorm.DB
)

func TestMain(m *testing.M) {
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

	//objectStore = objectstore.NewFileObjectStore(testPath+"objects/", testPath+"thumbnails/")

	code := m.Run()

	//err = os.RemoveAll(testPath)
	//if err != nil {
	//	panic(err)
	//}

	os.Exit(code)
}

func CleanDB() {
	db.Exec("DELETE FROM documents")
}
