package database

import (
	"identifEye/entity"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open(sqlite.Open("identifeye.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	db.AutoMigrate(&entity.User{})
}

func Get() *gorm.DB {
	if db == nil {
		Init()
	}
	return db
}
