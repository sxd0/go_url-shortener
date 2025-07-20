package handler

import (
	"context"
	"errors"

	"github.com/sxd0/go_url-shortener/internal/auth/jwt"
	"github.com/sxd0/go_url-shortener/internal/auth/repository"
	"github.com/sxd0/go_url-shortener/internal/auth/service"
	"github.com/sxd0/go_url-shortener/proto/authpb"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	AuthService    *service.AuthService
	TokenGenerator *jwt.JWT
	UserRepo       *repository.UserRepository
}

func NewAuthHandler(as *service.AuthService, tg *jwt.JWT, ur *repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		AuthService:    as,
		TokenGenerator: tg,
		UserRepo:       ur,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	_, err := h.AuthService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		return nil, err
	}

	user, err := h.UserRepo.FindByEmail(req.Email)
	if err != nil || user == nil {
		return nil, errors.New("failed to fetch user after register")
	}

	return &authpb.RegisterResponse{
		UserId: uint64(user.ID),
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	email, err := h.AuthService.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	token, err := h.TokenGenerator.Create(jwt.JWTData{
		Email:  email,
		UserID: 0,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := h.TokenGenerator.Create(jwt.JWTData{
		Email:  email,
		UserID: 0,
	})
	if err != nil {
		return nil, err
	}

	return &authpb.LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (h *AuthHandler) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	valid, data := h.TokenGenerator.Parse(req.RefreshToken)
	if !valid {
		return nil, errors.New("invalid refresh token")
	}

	token, err := h.TokenGenerator.Create(jwt.JWTData{
		UserID: data.UserID,
		Email:  data.Email,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := h.TokenGenerator.Create(jwt.JWTData{
		UserID: data.UserID,
		Email:  data.Email,
	})
	if err != nil {
		return nil, err
	}

	return &authpb.RefreshResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (h *AuthHandler) VerifyToken(ctx context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	valid, data := h.TokenGenerator.Parse(req.AccessToken)
	if !valid {
		return &authpb.VerifyTokenResponse{Valid: false}, nil
	}
	return &authpb.VerifyTokenResponse{
		Valid:  true,
		UserId: uint64(data.UserID),
	}, nil
}

func (h *AuthHandler) GetUserByID(ctx context.Context, req *authpb.GetUserByIDRequest) (*authpb.GetUserByIDResponse, error) {
	return nil, errors.New("method GetUserByID not implemented")
}
