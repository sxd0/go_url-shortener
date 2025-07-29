package main

import (
	"log"
	"net/http"

	"github.com/sxd0/go_url-shortener/internal/gateway"
	"github.com/sxd0/go_url-shortener/internal/gateway/configs"
	"github.com/sxd0/go_url-shortener/internal/gateway/jwt"
	"github.com/sxd0/go_url-shortener/internal/gateway/logger"
	"github.com/sxd0/go_url-shortener/internal/gateway/redis"
	"github.com/sxd0/go_url-shortener/internal/gateway/service"
)

func main() {
	cfg := configs.LoadConfig()

	logger.InitFromEnv()
	defer logger.Sync()

	verifier := jwt.NewVerifier(cfg.PublicKey)

	authSvc, err := service.NewAuthService(cfg.AuthGRPCAddr)
	if err != nil {
		log.Fatalf("auth client: %v", err)
	}
	linkSvc, err := service.NewLinkService(cfg.LinkGRPCAddr)
	if err != nil {
		log.Fatalf("link client: %v", err)
	}
	statSvc, err := service.NewStatService(cfg.StatGRPCAddr)
	if err != nil {
		log.Fatalf("stat client: %v", err)
	}

	rdb := redis.New(cfg.RedisAddr)

	router := gateway.NewRouter(gateway.RedirectDeps{
		Verifier:   verifier,
		AuthClient: authSvc.Client(),
		LinkClient: linkSvc.Client(),
		StatClient: statSvc.Client(),
		Cache:      rdb,
	}, cfg)

	log.Printf("Gateway listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}
