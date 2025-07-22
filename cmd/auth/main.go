package main

import (
	"log"
	"net"

	"github.com/sxd0/go_url-shortener/internal/auth"
	"github.com/sxd0/go_url-shortener/internal/auth/configs"
	"github.com/sxd0/go_url-shortener/internal/auth/handler"
	"github.com/sxd0/go_url-shortener/internal/auth/jwt"
	"github.com/sxd0/go_url-shortener/internal/auth/repository"
	"github.com/sxd0/go_url-shortener/internal/auth/service"
	"github.com/sxd0/go_url-shortener/proto/authpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := configs.LoadConfig()

	dbConn := auth.NewDb(cfg)
	userRepo := repository.NewUserRepository(dbConn)
	authService := service.NewAuthService(userRepo)
	tokenGenerator := jwt.NewJWT(cfg.Auth.Secret)

	authHandler := handler.NewAuthHandler(authService, tokenGenerator, userRepo)

	lis, err := net.Listen("tcp", ":"+cfg.App.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, authHandler)
	reflection.Register(server)

	log.Printf("Auth gRPC server listening on :%s", cfg.App.Port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
