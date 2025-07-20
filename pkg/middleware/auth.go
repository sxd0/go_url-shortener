package middleware

import (
	"context"
	"net/http"

	"github.com/sxd0/go_url-shortener/configs"
	"github.com/sxd0/go_url-shortener/pkg/jwt"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func writeUnauthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuthed(config *configs.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				writeUnauthed(w)
				return
			}

			isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(cookie.Value)
			if !isValid {
				writeUnauthed(w)
				return
			}

			ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
