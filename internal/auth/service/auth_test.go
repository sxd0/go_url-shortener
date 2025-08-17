package service

import (
	"errors"
	"github.com/sxd0/go_url-shortener/internal/auth/logger"
	"go.uber.org/zap"
	"testing"

	"github.com/sxd0/go_url-shortener/internal/auth/di"
	"github.com/sxd0/go_url-shortener/internal/auth/model"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	users     map[string]*model.User
	findErr   error
	createErr error
}

func newMockRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*model.User)}
}

func (m *mockUserRepo) Create(user *model.User) (*model.User, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	m.users[user.Email] = user
	return user, nil
}

func (m *mockUserRepo) FindByEmail(email string) (*model.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	u, ok := m.users[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (m *mockUserRepo) FindByID(id uint) (*model.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

var _ di.IUserRepository = (*mockUserRepo)(nil)

func TestAuthServiceLogin(t *testing.T) {
	repo := newMockRepo()
	hashed, _ := bcryptHash("password")
	repo.Create(&model.User{Email: "user@example.com", Password: hashed})
	svc := NewAuthService(repo)

	email, err := svc.Login("user@example.com", "password")
	if err != nil || email != "user@example.com" {
		t.Fatalf("expected successful login, got %v, %v", email, err)
	}

	if _, err := svc.Login("missing@example.com", "password"); err == nil {
		t.Fatalf("expected error for missing user")
	}

	if _, err := svc.Login("user@example.com", "wrong"); err == nil {
		t.Fatalf("expected error for wrong password")
	}
}

func TestAuthServiceRegister(t *testing.T) {
	repo := newMockRepo()
	logger.Log = zap.NewNop()
	svc := NewAuthService(repo)

	email, err := svc.Register("new@example.com", "pass", "name")
	if err != nil || email != "new@example.com" {
		t.Fatalf("expected successful register, got %v, %v", email, err)
	}

	if _, err := svc.Register("new@example.com", "pass", "name"); err == nil {
		t.Fatalf("expected error for existing user")
	}
}

func bcryptHash(pwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(b), err
}

func TestVerifyRefreshToken(t *testing.T) {
	svc := NewAuthService(newMockRepo())
	if err := svc.VerifyRefreshToken(1, 2); err == nil {
		t.Fatalf("expected permission error")
	}
	if err := svc.VerifyRefreshToken(1, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}