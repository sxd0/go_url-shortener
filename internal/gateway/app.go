package gateway

import (
	"net/http"
	"time"
)

func App(deps Deps) *http.Server {
	return &http.Server{
		Addr:              ":8080",
		Handler:           NewRouter(deps),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       10 * time.Second,
	}
}
