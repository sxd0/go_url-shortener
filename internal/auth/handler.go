package auth

import (
	"go/test-http/configs"
	"go/test-http/pkg/jwt"
	"go/test-http/pkg/req"
	"go/test-http/pkg/res"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

type AuthHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(r chi.Router, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}

	r.Group(func(r chi.Router) {
		r.Post("/auth/login", handler.Login())
		r.Post("/auth/register", handler.Register())
		r.Post("/auth/refresh", handler.Refresh())
	})
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}

		email, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := handler.UserRepository.FindByEmail(email)
		if err != nil || user == nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		accessToken, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Email:  user.Email,
			UserID: user.ID,
			Exp:    15 * time.Minute,
		})
		if err != nil {
			http.Error(w, "failed to generate access token", http.StatusInternalServerError)
			return
		}

		refreshToken, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Email:     user.Email,
			UserID:    user.ID,
			IsRefresh: true,
			Exp:       7 * 24 * time.Hour,
		})
		if err != nil {
			http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  time.Now().Add(15 * time.Minute),
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/auth/refresh",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})

		res.Json(w, map[string]string{"message": "login successful"}, http.StatusOK)
	}
}

func (handler *AuthHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		isValid, data := jwt.NewJWT(handler.Config.Auth.Secret).Parse(cookie.Value)
		if !isValid || !data.IsRefresh {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}

		accessToken, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Email:  data.Email,
			UserID: data.UserID,
			Exp:    15 * time.Minute,
		})
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  time.Now().Add(15 * time.Minute),
		})

		res.Json(w, map[string]string{"message": "token refreshed"}, http.StatusOK)
	}
}


func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}

		email, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := handler.UserRepository.FindByEmail(email)
		if err != nil || user == nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		accessToken, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Email:  user.Email,
			UserID: user.ID,
			Exp:    15 * time.Minute,
		})
		if err != nil {
			http.Error(w, "failed to generate access token", http.StatusInternalServerError)
			return
		}

		refreshToken, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Email:     user.Email,
			UserID:    user.ID,
			IsRefresh: true,
			Exp:       7 * 24 * time.Hour,
		})
		if err != nil {
			http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  time.Now().Add(15 * time.Minute),
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/auth/refresh",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})

		res.Json(w, map[string]string{"message": "registration successful"}, http.StatusCreated)
	}
}

