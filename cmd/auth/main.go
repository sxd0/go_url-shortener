package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sxd0/go_url-shortener/internal/auth"
	"github.com/sxd0/go_url-shortener/internal/auth/configs"
	"github.com/sxd0/go_url-shortener/internal/auth/handler"
	"github.com/sxd0/go_url-shortener/internal/auth/jwt"
	"github.com/sxd0/go_url-shortener/internal/auth/repository"
	"github.com/sxd0/go_url-shortener/internal/auth/server"
	"github.com/sxd0/go_url-shortener/internal/auth/service"
	"github.com/sxd0/go_url-shortener/proto/gen/go/authpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := configs.LoadConfig()

	dbConn := auth.NewDb(cfg)
	userRepo := repository.NewUserRepository(dbConn)
	authService := service.NewAuthService(userRepo)

	privKey, err := jwt.LoadRSAPrivateKey(cfg.Auth.PrivateKeyPath)
	if err != nil {
		log.Fatalf("failed to load private key: %v", err)
	}

	pubKey, err := jwt.LoadRSAPublicKey(cfg.Auth.PublicKeyPath)
	if err != nil {
		log.Fatalf("failed to load public key: %v", err)
	}

	tokenGenerator := jwt.NewJWT(privKey, pubKey)

	authHandler := handler.NewAuthHandler(authService, tokenGenerator, userRepo)

	lis, err := net.Listen("tcp", ":"+cfg.App.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := server.NewGRPCServerWithMiddleware(tokenGenerator)
	authpb.RegisterAuthServiceServer(grpcServer, authHandler)
	reflection.Register(grpcServer)

	go func() {
		log.Printf("Auth gRPC server listening on :%s", cfg.App.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down gRPC server...")

	gracefulStop(grpcServer)
}

func gracefulStop(server *grpc.Server) {
	done := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("gRPC server stopped gracefully")
	case <-time.After(10 * time.Second):
		log.Println("Timeout — forcing server shutdown")
		server.Stop()
	}
}
