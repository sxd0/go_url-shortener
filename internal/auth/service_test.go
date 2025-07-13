package auth_test

import (
	"go/test-http/internal/auth"
	"go/test-http/internal/user"
	"testing"
)

type MockUserRepository struct{}

func (repo *MockUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: "user@example.com",
	}, nil
}

func (repo *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}

func TestRegisterSuccess(t *testing.T) {
	const initialEmail = "user@example.com"
	authService := auth.NewAuthService(&MockUserRepository{})
	email, err :=authService.Register(initialEmail, "123", "Anton")
	if err != nil {
		t.Fatal(err)
	}
	if email != initialEmail {
		t.Fatalf("Expected %s got %s", "user@example.com", email)
	}
}
