package main

import (
	"fmt"
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
	"github.com/sxd0/go_url-shortener/internal/link/configs"
	"github.com/sxd0/go_url-shortener/internal/link/db"
	"github.com/sxd0/go_url-shortener/internal/link/handler"
	"github.com/sxd0/go_url-shortener/internal/link/logger"
	"github.com/sxd0/go_url-shortener/internal/link/repository"
	"github.com/sxd0/go_url-shortener/internal/link/server"
	"github.com/sxd0/go_url-shortener/internal/link/service"
	"google.golang.org/grpc"
	healthgrpc "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	config := configs.LoadConfig()

	logger.InitFromEnv()
	defer logger.Sync()

	db := db.New(config.Db.GetDSN())

	repo := repository.NewLinkRepository(db)
	srv := service.NewLinkService(repo)

	h := handler.NewLinkHandler(srv)

	grpcServer := server.NewGRPCServerWithMiddleware()
	handler.RegisterLinkHandler(grpcServer, h)

	reflection.Register(grpcServer)

	grpc_prom.Register(grpcServer)

	healthServer := healthgrpc.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.App.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Printf("Link gRPC server listening on :%s", config.App.Port)
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
		log.Println("Link HTTP health on :9101")
		http.ListenAndServe(":9102", r)
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
