package di

import "github.com/sxd0/go_url-shortener/internal/auth/model"

type IUserRepository interface {
	Create(user *model.User) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
}
