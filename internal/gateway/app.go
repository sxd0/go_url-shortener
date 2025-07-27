package gateway

import (
	"net/http"
	"time"

	"github.com/sxd0/go_url-shortener/internal/gateway/configs"
)

func App(deps Deps, cfg *configs.Config) *http.Server {
	return &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           NewRouter(deps, cfg),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       10 * time.Second,
	}
}
