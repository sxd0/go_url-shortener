package main

import (
	"log"
	"net/http"

	"github.com/sxd0/go_url-shortener/internal/gateway"
	"github.com/sxd0/go_url-shortener/internal/gateway/configs"
	"github.com/sxd0/go_url-shortener/internal/gateway/jwt"
	"github.com/sxd0/go_url-shortener/internal/gateway/service"
)

func main() {
	cfg := configs.LoadConfig()

	// JWT Verifier
	verifier := jwt.NewVerifier(cfg.PublicKey)

	// gRPC clients
	authService, err := service.NewAuthService(cfg.AuthGRPCAddr)
	if err != nil {
		log.Fatal("cannot init AuthService:", err)
	}

	linkService, err := service.NewLinkService(cfg.LinkGRPCAddr)
	if err != nil {
		log.Fatal("cannot init LinkService:", err)
	}

	statService, err := service.NewStatService(cfg.StatGRPCAddr)
	if err != nil {
		log.Fatal("cannot init StatService:", err)
	}

	// Gateway Router
	router := gateway.NewRouter(gateway.Deps{
		Verifier:   verifier,
		AuthClient: authService.Client(),
		LinkClient: linkService.Client(),
		StatClient: statService.Client(),
	})

	// Run HTTP server
	log.Println("Gateway listening on port:", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal("server failed:", err)
	}
}
