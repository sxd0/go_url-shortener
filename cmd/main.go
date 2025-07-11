package main

import (
	"fmt"
	"go/test-http/configs"
	"go/test-http/internal/auth"
	"go/test-http/internal/link"
	"go/test-http/pkg/db"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()

	_ = db.NewDb(conf)

	router := http.NewServeMux()
	// Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server is listening on port :8081")
	server.ListenAndServe()
}
