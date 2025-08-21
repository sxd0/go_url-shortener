package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	grpc_prom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sxd0/go_url-shortener/internal/stat/configs"
	"github.com/sxd0/go_url-shortener/internal/stat/db"
	"github.com/sxd0/go_url-shortener/internal/stat/handler"
	"github.com/sxd0/go_url-shortener/internal/stat/logger"
	"github.com/sxd0/go_url-shortener/internal/stat/repository"
	"github.com/sxd0/go_url-shortener/internal/stat/server"
	"github.com/sxd0/go_url-shortener/internal/stat/service"
	"github.com/sxd0/go_url-shortener/pkg/kafka"
	healthgrpc "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/sxd0/go_url-shortener/proto/gen/go/statpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := configs.LoadConfig()

	logger.InitFromEnv()
	defer logger.Sync()

	db := db.New(cfg.Db.GetDSN())

	kafkaAddr := os.Getenv("KAFKA_ADDR")
	if kafkaAddr == "" {
		kafkaAddr = "kafka:9092"
	}
	subscriber := kafka.NewSubscriber([]string{kafkaAddr}, "link.events", "stat-service")
	defer subscriber.Close()

	statRepo := repository.NewStatRepository(db)
	statService := service.NewStatService(statRepo, subscriber)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := statService.Start(ctx); err != nil {
			log.Fatalf("stat consumer: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":"+cfg.App.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := server.NewGRPCServerWithMiddleware()

	statpb.RegisterStatServiceServer(grpcServer, handler.NewStatGRPCHandler(statRepo))

	reflection.Register(grpcServer)

	grpc_prom.Register(grpcServer)

	healthServer := healthgrpc.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	go func() {
		log.Printf("Stat gRPC server listening on :%s", cfg.App.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	go func() {
		r := chi.NewRouter()
		r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
		r.Handle("/metrics", promhttp.Handler())
		log.Println("Stat HTTP health on :9103")
		http.ListenAndServe(":9103", r)
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
