package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/sxd0/go_url-shortener/internal/gateway/middleware"
	"github.com/sxd0/go_url-shortener/proto/gen/go/statpb"
	"google.golang.org/grpc/metadata"
)

type StatHandler struct {
	client statpb.StatServiceClient
}

func NewStatHandler(deps Deps) *StatHandler {
	return &StatHandler{client: deps.StatClient}
}

func (h *StatHandler) GetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, err := middleware.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		by := r.URL.Query().Get("by")
		if by == "" {
			by = "day"
		}
		if from == "" || to == "" {
			now := time.Now().UTC()
			to = now.Format("2006-01-02")
			from = now.AddDate(0, 0, -30).Format("2006-01-02")
		}

		ctx := middleware.AttachCommonMD(r.Context(), r)
		ctx = metadata.AppendToOutgoingContext(ctx, "x-user-id", strconv.FormatUint(uint64(uid), 10))

		req := &statpb.GetStatsRequest{
			From: from,
			To:   to,
			By:   by,
		}

		resp, err := h.client.GetStats(ctx, req)
		if err != nil {
			middleware.WriteGRPCError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		stats := resp.Stats
		if stats == nil {
			stats = []*statpb.Stat{}
		}
		_ = json.NewEncoder(w).Encode(stats)
	}
}
