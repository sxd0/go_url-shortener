package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
