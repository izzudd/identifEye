package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
	Email    string
	Name     string
	Key      string `gorm:"unique"`
}
