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

	// AUTH
	r.Route("/auth", func(r chi.Router) {
		authHandler := handler.NewAuthHandler(deps.AuthClient)

		r.Post("/register", authHandler.Register())
		r.Post("/login", authHandler.Login())
		r.Post("/refresh", authHandler.Refresh())

		r.Post("/validate", authHandler.Validate())
		r.With(middleware.JWTMiddleware(deps.Verifier)).
			Get("/user/{id}", authHandler.GetUserByID())
	})

	// LINKS
	r.Route("/link", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware(deps.Verifier))
		linkHandler := handler.NewLinkHandler(deps.LinkClient)

		r.Post("/", linkHandler.Create())
		// r.Get("/", linkHandler.GetAll())

		// r.Get("/{hash}", linkHandler.GetByHash())

		r.Patch("/", linkHandler.Update())

		r.Delete("/{id}", linkHandler.Delete())
		r.Delete("/hash/{hash}", linkHandler.DeleteByHash())
	})

	// REDIRECT
	r.Get("/r/{hash}", handler.RedirectHandler(deps.LinkClient, deps.StatClient))

	// STATS
	r.Route("/stat", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware(deps.Verifier))

		statHandler := handler.NewStatHandler(handler.Deps{
			StatClient: deps.StatClient,
		})

		r.Get("/", statHandler.GetStats())
	})

	return r
}
