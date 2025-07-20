package di

import "github.com/sxd0/go_url-shortener/internal/auth/model"

type IStatRepository interface {
	AddClick(linkId uint)
}

type IUserRepository interface {
	Create(user *model.User) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
}
