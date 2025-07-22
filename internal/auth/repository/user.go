package repository

import (
	"github.com/sxd0/go_url-shortener/internal/auth/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	Database *gorm.DB
}

func NewUserRepository(database *gorm.DB) *UserRepository {
	return &UserRepository{
		Database: database,
	}
}

func (repo *UserRepository) Create(user *model.User) (*model.User, error) {
	result := repo.Database.Create(user)
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

func (repo *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	result := repo.Database.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
