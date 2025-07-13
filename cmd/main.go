package main

import (
	"fmt"
	"go/test-http/configs"
	"go/test-http/internal/auth"
	"go/test-http/internal/link"
	"go/test-http/pkg/db"
	"go/test-http/pkg/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()

	db := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	linkRepository := link.NewLinkRepository(db)

	// Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: middleware.CORS(middleware.Logging(router)),
	}

	fmt.Println("Server is listening on port :8081")
	server.ListenAndServe()
}
