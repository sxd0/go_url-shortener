package service

import (
	"fmt"

	"github.com/sxd0/go_url-shortener/proto/gen/go/authpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthService struct {
	client authpb.AuthServiceClient
}

func (s *AuthService) Client() authpb.AuthServiceClient {
	return s.client
}

func NewAuthService(addr string) (*AuthService, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AuthService: %w", err)
	}

	client := authpb.NewAuthServiceClient(conn)
	return &AuthService{client: client}, nil
}
