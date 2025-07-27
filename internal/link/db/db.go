package db

import (
	"log"
	"time"

	"github.com/sxd0/go_url-shortener/internal/link/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("link db connect: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("link db sql: %v", err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	_ = model.Link{}
	return db
}
