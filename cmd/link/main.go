package main

import (
	"fmt"
	"log"
	"net"

	"github.com/sxd0/go_url-shortener/internal/link"
	"github.com/sxd0/go_url-shortener/internal/link/configs"
	"github.com/sxd0/go_url-shortener/internal/link/handler"
	"github.com/sxd0/go_url-shortener/internal/link/repository"
	"github.com/sxd0/go_url-shortener/internal/link/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config := configs.LoadConfig()
	db := link.NewDb(config)

	repo := repository.NewLinkRepository(db)
	srv := service.NewLinkService(repo)

	grpcServer := grpc.NewServer()
	handler.RegisterLinkHandler(grpcServer, srv)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.App.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Link Service started on port %s", config.App.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
