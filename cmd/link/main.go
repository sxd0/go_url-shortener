package main

import (
	"fmt"
	"log"
	"net"

	"github.com/sxd0/go_url-shortener/internal/link/configs"
	"github.com/sxd0/go_url-shortener/internal/link/handler"
	"github.com/sxd0/go_url-shortener/internal/link/repository"
	"github.com/sxd0/go_url-shortener/internal/link/service"

	"google.golang.org/grpc"
)

func main() {
	config := configs.LoadConfig()
	db := configs.NewDb(config)

	linkRepo := repository.NewLinkRepository(db)
	linkService := service.NewLinkService(linkRepo)

	grpcServer := grpc.NewServer()
	handler.RegisterLinkHandler(grpcServer, linkService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Link Service started on port %s", config.GrpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
