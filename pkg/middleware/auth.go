package middleware

import (
	"context"
	"go/test-http/configs"
	"go/test-http/pkg/jwt"
	"net/http"
	"strings"
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
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeUnauthed(w)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(token)
			if !isValid {
				writeUnauthed(w)
				return
			}

			ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

