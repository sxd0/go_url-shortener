package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sxd0/go_url-shortener/internal/gateway/configs"
	"github.com/sxd0/go_url-shortener/internal/gateway/jwt"
	"github.com/sxd0/go_url-shortener/internal/gateway/middleware"
)

func main() {
	cfg := configs.LoadConfig()
	verifier := jwt.NewVerifier(cfg.PublicKey)

	r := chi.NewRouter()
	r.Use(middleware.JWTMiddleware(verifier))

	r.Get("/stat", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(uint)
		w.Write([]byte("Stat requested by user ID: " + string(rune(userID))))
	})

	log.Println("Gateway listening on port:", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal("server failed:", err)
	}
}
