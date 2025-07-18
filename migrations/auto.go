package main

import (
	"fmt"
	"go/test-http/configs"
	"go/test-http/internal/link"
	"go/test-http/internal/stat"
	"go/test-http/internal/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	conf := configs.LoadConfig()
	fmt.Println("DSN:", conf.Db.GetDSN())

	db, err := gorm.Open(postgres.Open(conf.Db.GetDSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Migrations is access")

	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&link.Link{})
	db.AutoMigrate(&stat.Stat{})
}
