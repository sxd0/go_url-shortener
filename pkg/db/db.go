package db

import (
	"github.com/sxd0/go_url-shortener/configs"
	"gorm.io/gorm"
)

func NewDb(cfg *configs.Config) *gorm.DB {
	return New(cfg.Db.GetDSN())
}
