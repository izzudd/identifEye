package database

import (
	"identifEye/entity"
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var once sync.Once
var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open(sqlite.Open("identifeye.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	db.AutoMigrate(&entity.User{})
	db.Debug()
}

func Get() *gorm.DB {
	once.Do(Init)
	return db
}
