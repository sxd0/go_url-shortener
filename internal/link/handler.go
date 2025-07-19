package link

import (
	"fmt"
	"go/test-http/configs"
	"go/test-http/internal/user"
	"go/test-http/pkg/event"
	"go/test-http/pkg/logger"
	"go/test-http/pkg/middleware"
	"go/test-http/pkg/req"
	"go/test-http/pkg/res"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
	UserRepository *user.UserRepository
	Config         *configs.Config
	EventBus       *event.EventBus
}

type LinkHandler struct {
	LinkRepository *LinkRepository
	UserRepository *user.UserRepository
	EventBus       *event.EventBus
}

func NewLinkHandler(r chi.Router, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
		UserRepository: deps.UserRepository,
		EventBus:       deps.EventBus,
	}
	// (middleware IsAuthed)
	r.Group(func(r chi.Router) {
		r.Use(middleware.IsAuthed(deps.Config))

		r.Post("/link", handler.Create())
		r.Get("/link", handler.GetAll())
		r.Patch("/link/{id}", handler.Update())
		r.Delete("/link/{id}", handler.Delete())
	})

	// Redirect
	r.Get("/{hash}", handler.GoTo())
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LinkCreateRequest](&w, r)
		if err != nil {
			return
		}

		email, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := handler.UserRepository.FindByEmail(email)
		if err != nil || user == nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		linkObj, err := NewLink(body.Url, func(hash string) bool {
			existing, _ := handler.LinkRepository.GetByHash(hash)
			return existing != nil
		})
		if err != nil {
			logger.Log.Error("failed to generate unique hash", zap.Error(err))
			http.Error(w, "could not create link", http.StatusInternalServerError)
			return
		}

		linkObj.UserID = user.ID

		createdLink, err := handler.LinkRepository.Create(linkObj)
		if err != nil {
			logger.Log.Error("failed to save link", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, createdLink, http.StatusCreated)
	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if ok {
			fmt.Println(email)
		}
		body, err := req.HandleBody[LinkUpdateRequest](&w, r)
		if err != nil {
			return
		}
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, link, 201)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = handler.LinkRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = handler.LinkRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, nil, 200)
	}
}

func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		link, err := handler.LinkRepository.GetByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		go handler.EventBus.Publish(event.Event{
			Type: event.EventLinkVisited,
			Data: link.ID,
		})
		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
	}
}

func (h *LinkHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := h.UserRepository.FindByEmail(email)
		if err != nil || user == nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit <= 0 {
			limit = 10
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil || offset < 0 {
			offset = 0
		}

		links, err := h.LinkRepository.GetAllByUserID(user.ID, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, links, http.StatusOK)
	}
}
