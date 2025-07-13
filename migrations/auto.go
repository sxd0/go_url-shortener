package main

import (
	"fmt"
	"go/test-http/internal/link"
	"go/test-http/internal/stat"
	"go/test-http/internal/user"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	fmt.Println("DSN:", os.Getenv("DSN"))
	fmt.Println("Migrations is access")
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&link.Link{}, &user.User{}, &stat.Stat{})
}
