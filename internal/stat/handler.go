package stat

import (
	"go/test-http/configs"
	"go/test-http/internal/user"
	"go/test-http/pkg/middleware"
	"go/test-http/pkg/res"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatHandlerDeps struct {
	StatRepository *StatRepository
	UserRepository *user.UserRepository
	Config         *configs.Config
}

type StatHandler struct {
	StatRepository *StatRepository
	UserRepository *user.UserRepository
}

func NewStatHandler(r chi.Router, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
		UserRepository: deps.UserRepository,
	}

	r.Group(func(r chi.Router) {
		r.Use(middleware.IsAuthed(deps.Config))

		r.Get("/stat", handler.GetStat())
	})
}

func (h *StatHandler) GetStat() http.HandlerFunc {
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

		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")
		by := r.URL.Query().Get("by")

		from, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			http.Error(w, "invalid 'from' date", http.StatusBadRequest)
			return
		}

		to, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			http.Error(w, "invalid 'to' date", http.StatusBadRequest)
			return
		}

		if by != "day" && by != "month" && by != "year" {
			http.Error(w, "invalid 'by' value", http.StatusBadRequest)
			return
		}

		stats, err := h.StatRepository.GetByUserID(user.ID, from, to, by)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, stats, http.StatusOK)
	}
}
