package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// panic(err)
		log.Fatalf("failed to connect database: %v", err)
	}
	return db
}
