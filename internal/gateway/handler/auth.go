package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	httpx "github.com/sxd0/go_url-shortener/internal/gateway/http"
	"github.com/sxd0/go_url-shortener/internal/gateway/middleware"
	"github.com/sxd0/go_url-shortener/proto/gen/go/authpb"
)

type AuthHandler struct {
	Client authpb.AuthServiceClient
}

func NewAuthHandler(client authpb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{Client: client}
}

type registerReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

func (h *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := httpx.Decode[registerReq](r)
		if err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := httpx.Validate(body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := middleware.AttachCommonMD(r.Context(), r)
		resp, err := h.Client.Register(ctx, &authpb.RegisterRequest{
			Email:    body.Email,
			Password: body.Password,
			Name:     body.Name,
		})
		if err != nil {
			middleware.WriteGRPCError(w, err)
			return
		}

		httpx.JSON(w, resp, http.StatusCreated)
	}
}

type loginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := httpx.Decode[loginReq](r)
		if err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := httpx.Validate(body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := middleware.AttachCommonMD(r.Context(), r)
		resp, err := h.Client.Login(ctx, &authpb.LoginRequest{
			Email:    body.Email,
			Password: body.Password,
		})
		if err != nil {
			middleware.WriteGRPCError(w, err)
			return
		}
		httpx.JSON(w, resp, http.StatusOK)
	}
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *AuthHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := httpx.Decode[refreshReq](r)
		if err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := httpx.Validate(body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := middleware.AttachCommonMD(r.Context(), r)
		resp, err := h.Client.Refresh(ctx, &authpb.RefreshRequest{
			RefreshToken: body.RefreshToken,
		})
		if err != nil {
			middleware.WriteGRPCError(w, err)
			return
		}
		httpx.JSON(w, resp, http.StatusOK)
	}
}

func (h *AuthHandler) Validate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")

		ctx := middleware.AttachCommonMD(r.Context(), r)
		resp, err := h.Client.VerifyToken(ctx, &authpb.VerifyTokenRequest{
			AccessToken: token,
		})
		if err != nil {
			middleware.WriteGRPCError(w, err)
			return
		}
		httpx.JSON(w, resp, http.StatusOK)
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

		ctx := middleware.AttachCommonMD(r.Context(), r)
		resp, err := h.Client.GetUserByID(ctx, &authpb.GetUserByIDRequest{
			UserId: id,
		})
		if err != nil {
			middleware.WriteGRPCError(w, err)
			return
		}
		httpx.JSON(w, resp, http.StatusOK)
	}
}
