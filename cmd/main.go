package main

import (
	"fmt"
	"go/test-http/configs"
	"go/test-http/internal/auth"
	"go/test-http/internal/link"
	"go/test-http/internal/user"
	"go/test-http/pkg/db"
	"go/test-http/pkg/event"
	"go/test-http/pkg/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()

	db := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	// Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	// statRepository := stat.NewStatRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)

	// Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		Config:         conf,
		EventBus:       eventBus,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Println("Server is listening on port :8081")
	server.ListenAndServe()
}
