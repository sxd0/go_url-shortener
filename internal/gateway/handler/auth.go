package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
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

func (h *AuthHandler) Validate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		resp, err := h.Client.VerifyToken(r.Context(), &authpb.VerifyTokenRequest{
			AccessToken: token,
		})
		if err != nil {
			http.Error(w, "failed to verify token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		type out struct {
			Valid  bool   `json:"valid"`
			UserID uint64 `json:"user_id,omitempty"`
		}
		if !resp.Valid {
			res.Json(w, out{Valid: false}, http.StatusOK)
			return
		}
		res.Json(w, out{Valid: true, UserID: resp.UserId}, http.StatusOK)
	}
}

func (h *AuthHandler) GetUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil || id == 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		resp, err := h.Client.GetUserByID(r.Context(), &authpb.GetUserByIDRequest{
			UserId: id,
		})
		if err != nil {
			http.Error(w, "failed to get user: "+err.Error(), http.StatusNotFound)
			return
		}
		res.Json(w, resp, http.StatusOK)
	}
}
