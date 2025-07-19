package main

import (
	"fmt"
	"go/test-http/configs"
	"go/test-http/internal/auth"
	"go/test-http/internal/link"
	"go/test-http/internal/stat"
	"go/test-http/internal/user"
	"go/test-http/pkg/db"
	"go/test-http/pkg/event"
	"go/test-http/pkg/logger"
	"go/test-http/pkg/middleware"
	"net/http"

	"github.com/go-chi/chi"
)

func App() http.Handler {
	conf := configs.LoadConfig()

	db := db.NewDb(conf)
	r := chi.NewRouter()
	eventBus := event.NewEventBus()

	// Middlewares
	r.Use(middleware.CORS)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)

	// Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})

	// Handler
	auth.NewAuthHandler(r, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(r, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		UserRepository: userRepository,
		Config:         conf,
		EventBus:       eventBus,
	})
	stat.NewStatHandler(r, stat.StatHandlerDeps{
		StatRepository: statRepository,
		UserRepository: userRepository,
		Config:         conf,
	})

	go statService.AddClick()

	return r
}

func main() {
	app := App()

	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	fmt.Println("Server is listening on port :8081")

	logger.InitLogger()
	defer logger.SyncLogger()

	server.ListenAndServe()
}
