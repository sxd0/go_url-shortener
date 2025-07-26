package handler

import (
	"net/http"

	"github.com/sxd0/go_url-shortener/pkg/req"
	"github.com/sxd0/go_url-shortener/pkg/res"
	"github.com/sxd0/go_url-shortener/proto/gen/go/authpb"
)

type AuthHandler struct {
	Client authpb.AuthServiceClient
}

func NewAuthHandler(client authpb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{Client: client}
}

func (h *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[authpb.RegisterRequest](&w, r)
		if err != nil {
			return
		}

		resp, err := h.Client.Register(r.Context(), body)
		if err != nil {
			http.Error(w, "failed to register: "+err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, resp, http.StatusCreated)
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[authpb.LoginRequest](&w, r)
		if err != nil {
			return
		}

		resp, err := h.Client.Login(r.Context(), body)
		if err != nil {
			http.Error(w, "failed to login: "+err.Error(), http.StatusUnauthorized)
			return
		}

		res.Json(w, resp, http.StatusOK)
	}
}

func (h *AuthHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[authpb.RefreshRequest](&w, r)
		if err != nil {
			return
		}

		resp, err := h.Client.Refresh(r.Context(), body)
		if err != nil {
			http.Error(w, "failed to refresh token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		res.Json(w, resp, http.StatusOK)
	}
}
