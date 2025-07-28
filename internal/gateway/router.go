package gateway

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sxd0/go_url-shortener/internal/gateway/configs"
	"github.com/sxd0/go_url-shortener/internal/gateway/handler"
	"github.com/sxd0/go_url-shortener/internal/gateway/middleware"
	"github.com/sxd0/go_url-shortener/internal/gateway/openapi"
)

func NewRouter(deps Deps, cfg *configs.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestIDMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddlewareWithCfg(cfg))

	// Docs
	openapi.Mount(r)

	// AUTH
	r.Route("/auth", func(r chi.Router) {
		h := handler.NewAuthHandler(deps.AuthClient)

		r.Post("/register", h.Register())
		r.Post("/login", h.Login())
		r.Post("/refresh", h.Refresh())
		r.Post("/validate", h.Validate())

		r.With(middleware.JWTMiddleware(deps.Verifier)).
			Get("/user/{id}", h.GetUserByID())
	})

	// LINK
	r.Route("/link", func(r chi.Router) {
		h := handler.NewLinkHandler(deps.LinkClient)

		r.Use(middleware.JWTMiddleware(deps.Verifier))

		r.Post("/", h.Create())
		r.Get("/", h.List())
		r.Get("/{hash}", h.Get())
		r.Patch("/", h.Update())
		r.Delete("/{id}", h.Delete())
		r.Delete("/hash/{hash}", h.DeleteByHash())
	})

	// Redirect
	r.Get("/r/{hash}", handler.RedirectHandler(handler.Deps{
		LinkClient: deps.LinkClient,
		StatClient: deps.StatClient,
		Verifier:   deps.Verifier,
	}))

	// STATS
	r.Route("/stat", func(r chi.Router) {
		h := handler.NewStatHandler(handler.Deps{
			StatClient: deps.StatClient,
		})
		r.Use(middleware.JWTMiddleware(deps.Verifier))
		r.Get("/", h.GetStats())
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return r
}
