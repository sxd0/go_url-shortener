package handler

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sxd0/go_url-shortener/internal/gateway/jwt"
	"github.com/sxd0/go_url-shortener/internal/gateway/middleware"
	"github.com/sxd0/go_url-shortener/internal/gateway/redis"
	"github.com/sxd0/go_url-shortener/pkg/kafka"
	"github.com/sxd0/go_url-shortener/proto/gen/go/authpb"
	"github.com/sxd0/go_url-shortener/proto/gen/go/linkpb"
	"github.com/sxd0/go_url-shortener/proto/gen/go/statpb"
)

type RedirectDeps struct {
	AuthClient     authpb.AuthServiceClient
	LinkClient     linkpb.LinkServiceClient
	StatClient     statpb.StatServiceClient
	Verifier       *jwt.Verifier
	Cache          *redis.Client
	CacheTTL       time.Duration
	KafkaPublisher *kafka.Publisher
}

func RedirectHandler(deps RedirectDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		if hash == "" {
			http.NotFound(w, r)
			return
		}

		if dest, ok := deps.Cache.GetString("link:" + hash); ok {
			http.Redirect(w, r, dest, http.StatusTemporaryRedirect)
			return
		}

		grpcCtx := middleware.AttachCommonMD(r.Context(), r)

		linkResp, err := deps.LinkClient.GetLinkByHash(grpcCtx, &linkpb.GetLinkByHashRequest{Hash: hash})
		if err != nil || linkResp.GetLink() == nil {
			http.NotFound(w, r)
			return
		}

		dest := linkResp.GetLink().GetUrl()
		u, err := url.ParseRequestURI(dest)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}

		if deps.KafkaPublisher != nil {
			_ = deps.KafkaPublisher.Publish(r.Context(), kafka.Event{
				Kind:   kafka.LinkVisitedKind,
				LinkID: uint(linkResp.GetLink().GetId()),
				UserID: uint(linkResp.GetLink().GetUserId()),
				Ts:     time.Now().UTC(),
			})
		} else {
			http.Error(w, "Kafka publisher is nil – click event not sent", http.StatusBadRequest)
		}

		deps.Cache.SetString("link:"+hash, dest, deps.CacheTTL)
		http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
	}
}
