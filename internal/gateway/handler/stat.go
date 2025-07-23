package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sxd0/go_url-shortener/internal/gateway/middleware"
	"github.com/sxd0/go_url-shortener/internal/gateway/service"
)

type StatHandler struct {
	StatService *service.StatService
}

func NewStatHandler(statService *service.StatService) *StatHandler {
	return &StatHandler{StatService: statService}
}

func (h *StatHandler) GetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		by := r.URL.Query().Get("by")

		if from == "" || to == "" || by == "" {
			http.Error(w, "missing from/to/by parameters", http.StatusBadRequest)
			return
		}
		if _, err := time.Parse("2006-01-02", from); err != nil {
			http.Error(w, "invalid from date", http.StatusBadRequest)
			return
		}
		if _, err := time.Parse("2006-01-02", to); err != nil {
			http.Error(w, "invalid to date", http.StatusBadRequest)
			return
		}
		if by != "day" && by != "month" {
			http.Error(w, "invalid by parameter", http.StatusBadRequest)
			return
		}

		resp, err := h.StatService.GetStats(r.Context(), userID, from, to, by)
		if err != nil {
			http.Error(w, "failed to get stats: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp.Stats)
	}
}
