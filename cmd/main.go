package main

import (
	"fmt"
	"net/http"

	"github.com/sxd0/go_url-shortener/configs"
	"github.com/sxd0/go_url-shortener/internal/auth"
	"github.com/sxd0/go_url-shortener/internal/link"
	"github.com/sxd0/go_url-shortener/internal/stat"
	"github.com/sxd0/go_url-shortener/internal/auth/repository"
	"github.com/sxd0/go_url-shortener/pkg/db"
	"github.com/sxd0/go_url-shortener/pkg/event"
	"github.com/sxd0/go_url-shortener/pkg/logger"
	"github.com/sxd0/go_url-shortener/pkg/middleware"

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
	userRepository := repository.NewUserRepository(db)
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
