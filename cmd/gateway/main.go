package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sxd0/go_url-shortener/internal/gateway/configs"
	"github.com/sxd0/go_url-shortener/internal/gateway/handler"
	"github.com/sxd0/go_url-shortener/internal/gateway/jwt"
	"github.com/sxd0/go_url-shortener/internal/gateway/middleware"
	"github.com/sxd0/go_url-shortener/internal/gateway/service"
)

func main() {
	cfg := configs.LoadConfig()

	// JWT Verifier
	verifier := jwt.NewVerifier(cfg.PublicKey)

	// Services
	statService, err := service.NewStatService(cfg.StatGRPCAddr)
	if err != nil {
		log.Fatal("cannot init StatService:", err)
	}

	// Routers
	r := chi.NewRouter()
	r.Use(middleware.JWTMiddleware(verifier))

	// Handlers
	statHandler := handler.NewStatHandler(handler.Deps{
		StatClient: statService.Client(),
	})
	r.Get("/stat", statHandler.GetStats())

	// Server
	log.Println("Gateway listening on port:", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal("server failed:", err)
	}
}
