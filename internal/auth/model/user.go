package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"index"`
	Password string
	Name     string
}
