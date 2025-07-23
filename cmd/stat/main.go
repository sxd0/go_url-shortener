package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sxd0/go_url-shortener/internal/stat/configs"
	"github.com/sxd0/go_url-shortener/internal/stat/handler"
	"github.com/sxd0/go_url-shortener/internal/stat/repository"
	"github.com/sxd0/go_url-shortener/internal/stat/server"
	"github.com/sxd0/go_url-shortener/internal/stat/service"
	"github.com/sxd0/go_url-shortener/pkg/event"
	"github.com/sxd0/go_url-shortener/proto/gen/go/statpb"
	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load config
	cfg := configs.LoadConfig()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db := configs.NewDb(cfg)

	bus := event.NewEventBus()

	statRepo := repository.NewStatRepository(db)
	statService := service.NewStatService(&service.StatServiceDeps{
		EventBus:       bus,
		StatRepository: statRepo,
	})

	go statService.AddClick()

	lis, err := net.Listen("tcp", cfg.Server.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := server.NewGRPCServerWithMiddleware()

	statpb.RegisterStatServiceServer(grpcServer, handler.NewStatGRPCHandler(statRepo))

	reflection.Register(grpcServer)

	go func() {
		logger.Info("StatService listening on " + cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("failed to serve: " + err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("StatService shutting down...")
	grpcServer.GracefulStop()
}
