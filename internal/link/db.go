package link

import (
	"github.com/sxd0/go_url-shortener/internal/link/configs"
	pkgdb "github.com/sxd0/go_url-shortener/pkg/db"
	"gorm.io/gorm"
)

func NewDb(cfg *configs.Config) *gorm.DB {
	return pkgdb.New(cfg.Db.GetDSN())
}
