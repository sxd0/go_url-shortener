package authclient

import (
	"context"
	"time"

	"github.com/sxd0/go_url-shortener/proto/gen/go/authpb"
	"google.golang.org/grpc"
)

type AuthClient struct {
	client authpb.AuthServiceClient
}

func NewAuthClient(conn *grpc.ClientConn) *AuthClient {
	return &AuthClient{
		client: authpb.NewAuthServiceClient(conn),
	}
}

func (a *AuthClient) Login(ctx context.Context, in *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return a.client.Login(ctx, in)
}

func (a *AuthClient) Register(ctx context.Context, in *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return a.client.Register(ctx, in)
}

func (a *AuthClient) Refresh(ctx context.Context, in *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return a.client.Refresh(ctx, in)
}
