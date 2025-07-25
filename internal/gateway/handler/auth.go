package handler

import (
	"net/http"

	"github.com/sxd0/go_url-shortener/internal/gateway/authclient"
)

type AuthHandler struct {
	client *authclient.AuthClient
}

func NewAuthHandler(deps Deps) *AuthHandler {
	return &AuthHandler{
		client: deps.AuthClient,
	}
}

func (h *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (h *AuthHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
