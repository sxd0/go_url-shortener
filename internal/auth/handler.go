package auth

import (
	"encoding/json"
	"fmt"
	"go/test-http/configs"
	"go/test-http/pkg/res"
	"net/http"
	"regexp"
)

type AuthHandlerDeps struct {
	*configs.Config
}

type AuthHandler struct {
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Login())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Прочитать body
		var payload LoginRequest
		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			res.Json(w, err.Error(), 402)
			return
		}

		if payload.Email == "" {
			res.Json(w, "Email required", 402)
			return
		}
		reg, _ := regexp.Compile(`[A-Za-z0-9\._%+\-]+@[A-Za-z0-9\.\-]+\.[A-Za-z]{2,}`)
		if !reg.MatchString(payload.Email) {
			res.Json(w, "Wrong email", 402)
			return
		}
		if payload.Password == "" {
			res.Json(w, "Password required", 402)
			return
		}

		fmt.Println(payload)
		data := LoginResponse{
			Token: "123",
		}
		res.Json(w, data, 200)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("register")
	}
}
