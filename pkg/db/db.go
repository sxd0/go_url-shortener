package db

import (
	"github.com/sxd0/go_url-shortener/configs"

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
