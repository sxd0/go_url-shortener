package gateway

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sxd0/go_url-shortener/internal/gateway/handler"
	"github.com/sxd0/go_url-shortener/internal/gateway/middleware"
)

func NewRouter(deps Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)

	r.Route("/auth", func(r chi.Router) {
		authHandler := handler.NewAuthHandler(handler.Deps{
			AuthClient: deps.AuthClient,
		})

		r.Post("/login", authHandler.Login())
		r.Post("/register", authHandler.Register())
		r.Post("/refresh", authHandler.Refresh())
	})

	r.Route("/link", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware(deps.Verifier))

		linkHandler := handler.NewLinkHandler(handler.Deps{
			LinkClient: deps.LinkClient,
		})

		r.Get("/", linkHandler.List())
		r.Post("/", linkHandler.Create())
		r.Get("/{id}", linkHandler.Get())
		r.Patch("/{id}", linkHandler.Update())
		r.Delete("/{id}", linkHandler.Delete())
	})

	r.Route("/stat", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware(deps.Verifier))

		statHandler := handler.NewStatHandler(handler.Deps{
			StatClient: deps.StatClient,
		})

		r.Get("/", statHandler.GetStats())
	})

	return r
}
