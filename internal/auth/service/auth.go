package service

import (
	"github.com/sxd0/go_url-shortener/internal/auth/model"
	"github.com/sxd0/go_url-shortener/pkg/di"
	"github.com/sxd0/go_url-shortener/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository di.IUserRepository
}

func NewAuthService(userRepository di.IUserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) Login(email, password string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser == nil {
		return "", status.Error(codes.Unauthenticated, "wrong email or password")
	}
	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "wrong email or password")
	}
	return existedUser.Email, nil
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil {
		logger.Log.Warn("user already exists", zap.String("email", email))
		return "", status.Error(codes.AlreadyExists, "user exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Warn("wrong credentials", zap.String("email", email))
		return "", status.Error(codes.Internal, "failed to hash password")
	}
	user := &model.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}

	_, err = service.UserRepository.Create(user)
	if err != nil {
		return "", status.Error(codes.Internal, "failed to create user")
	}
	return user.Email, nil
}
