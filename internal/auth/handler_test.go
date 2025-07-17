package auth_test

import (
	"go/test-http/configs"
	"go/test-http/internal/auth"
	"go/test-http/internal/user"
	"go/test-http/pkg/db"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestLoginSuccess(t *testing.T) {
	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed init mock db")
		return
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))
	if err != nil {
		t.Fatal("Failed init gorm")
		return
	}
	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDb,
	})
	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepo),
	}
}
