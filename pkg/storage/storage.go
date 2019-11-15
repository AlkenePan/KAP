package storage

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func init() {
	db, _ := OpenDb("/tmp/test.db")
	db.AutoMigrate(&AppTable{}, &SourceTable{}, &ExecutableTable{})
	db.AutoMigrate(&CryptoTable{})
	db.AutoMigrate(&AlertTable{})
	db.AutoMigrate(&BuildTable{})
}

func OpenDb(DbAbsPath string) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", DbAbsPath)
	return db, err
}

func testStorage(db *gorm.DB) {
	db.AutoMigrate(&AppTable{})
	db.Create(&AppTable{Appid: uuid.New()})

	var app_row AppTable
	db.First(&app_row, 1)
	fmt.Println(app_row.Appid)
}
