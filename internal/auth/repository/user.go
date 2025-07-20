package repository

import (
	"github.com/sxd0/go_url-shortener/internal/auth/model"
	"github.com/sxd0/go_url-shortener/pkg/db"
)

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		Database: database,
	}
}

func (repo *UserRepository) Create(user *model.User) (*model.User, error) {
	result := repo.Database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := repo.Database.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
