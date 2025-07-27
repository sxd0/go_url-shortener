package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sxd0/go_url-shortener/internal/stat/configs"
	"github.com/sxd0/go_url-shortener/internal/stat/db"
	"github.com/sxd0/go_url-shortener/internal/stat/handler"
	"github.com/sxd0/go_url-shortener/internal/stat/logger"
	"github.com/sxd0/go_url-shortener/internal/stat/repository"
	"github.com/sxd0/go_url-shortener/internal/stat/server"
	"github.com/sxd0/go_url-shortener/internal/stat/service"

	"github.com/sxd0/go_url-shortener/internal/stat/event"

	"github.com/sxd0/go_url-shortener/proto/gen/go/statpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load config
	cfg := configs.LoadConfig()

	logger.InitFromEnv()
	defer logger.Sync()

	db := db.New(cfg.Db.GetDSN())

	bus := event.NewEventBus()

	statRepo := repository.NewStatRepository(db)
	statService := service.NewStatService(&service.StatServiceDeps{
		EventBus:       bus,
		StatRepository: statRepo,
	})

	go statService.AddClick()

	lis, err := net.Listen("tcp", ":"+cfg.App.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := server.NewGRPCServerWithMiddleware()

	statpb.RegisterStatServiceServer(grpcServer, handler.NewStatGRPCHandler(statRepo))

	reflection.Register(grpcServer)

	go func() {
		log.Printf("Stat gRPC server listening on :%s", cfg.App.Port)
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
