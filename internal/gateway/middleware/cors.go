package middleware

import (
	"net/http"
	"slices"

	"github.com/sxd0/go_url-shortener/internal/gateway/configs"
)

func CORSMiddlewareWithCfg(cfg *configs.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		allowed := cfg.AllowedOrigins
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "http://localhost:5173" || (origin != "" && slices.Contains(allowed, origin)) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Max-Age", "3600")

				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusOK)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
