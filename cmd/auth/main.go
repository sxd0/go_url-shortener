package main

import (
	"log"
	"net"
	"os"

	"github.com/sxd0/go_url-shortener/configs"
	"github.com/sxd0/go_url-shortener/internal/auth/handler"
	"github.com/sxd0/go_url-shortener/internal/auth/jwt"
	"github.com/sxd0/go_url-shortener/internal/auth/repository"
	"github.com/sxd0/go_url-shortener/internal/auth/service"
	"github.com/sxd0/go_url-shortener/pkg/db"
	"github.com/sxd0/go_url-shortener/proto/authpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := configs.LoadConfig()

	dbConn := db.NewDb(cfg)
	userRepo := repository.NewUserRepository(dbConn)
	authService := service.NewAuthService(userRepo)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not set")
	}
	tokenGenerator := jwt.NewJWT(secret)

	authHandler := handler.NewAuthHandler(authService, tokenGenerator, userRepo)

	port := os.Getenv("AUTH_GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, authHandler)
	reflection.Register(server)

	log.Printf("Auth gRPC server listening on :%s", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
