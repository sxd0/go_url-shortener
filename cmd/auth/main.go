package main

import (
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
	"github.com/sxd0/go_url-shortener/internal/auth/configs"
	"github.com/sxd0/go_url-shortener/internal/auth/db"
	"github.com/sxd0/go_url-shortener/internal/auth/handler"
	"github.com/sxd0/go_url-shortener/internal/auth/jwt"
	"github.com/sxd0/go_url-shortener/internal/auth/logger"
	"github.com/sxd0/go_url-shortener/internal/auth/repository"
	"github.com/sxd0/go_url-shortener/internal/auth/server"
	"github.com/sxd0/go_url-shortener/internal/auth/service"
	"github.com/sxd0/go_url-shortener/proto/gen/go/authpb"
	"google.golang.org/grpc"
	healthgrpc "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := configs.LoadConfig()

	logger.InitFromEnv()
	defer logger.Sync()

	dbConn := db.New(cfg.Db.GetDSN())
	if dbConn == nil {
		log.Fatal("nil *gorm.DB returned")
	}
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

	grpc_prom.Register(grpcServer)

	healthServer := healthgrpc.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	go func() {
		log.Printf("Auth gRPC server listening on :%s", cfg.App.Port)
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
		log.Println("Auth HTTP health on :9101")
		http.ListenAndServe(":9101", r)
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
