package db

import (
	"go/test-http/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDb(conf *configs.Config) *Db {
	db, err := gorm.Open(postgres.Open(conf.Db.GetDSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Db{db}
}
