package handler

import (
	"context"
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
	"golang.org/x/sync/singleflight"
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

var sf singleflight.Group

type redirData struct {
	dest    string
	linkID  uint
	ownerID uint
}

func RedirectHandler(deps RedirectDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		if hash == "" {
			http.NotFound(w, r)
			return
		}
		cacheKey := "cache:link:" + hash

		if deps.Cache != nil {
			ctx, cancel := context.WithTimeout(r.Context(), 120*time.Millisecond)
			defer cancel()

			if dest, ok, err := deps.Cache.GetString(ctx, cacheKey); err == nil && ok {
				http.Redirect(w, r, dest, http.StatusTemporaryRedirect)
				return
			}
		}

		v, err, _ := sf.Do(cacheKey, func() (interface{}, error) {
			grpcCtx := middleware.AttachCommonMD(r.Context(), r)
			linkResp, err := deps.LinkClient.GetLinkByHash(grpcCtx, &linkpb.GetLinkByHashRequest{Hash: hash})
			if err != nil || linkResp.GetLink() == nil {
				return nil, err
			}
			link := linkResp.GetLink()
			return redirData{
				dest:    link.GetUrl(),
				linkID:  uint(link.GetId()),
				ownerID: uint(link.GetUserId()),
			}, nil
		})
		if err != nil || v == nil {
			http.NotFound(w, r)
			return
		}
		data := v.(redirData)

		u, err := url.ParseRequestURI(data.dest)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}

		if deps.Cache != nil {
			ctx, cancel := context.WithTimeout(r.Context(), 120*time.Millisecond)
			_ = deps.Cache.SetString(ctx, cacheKey, data.dest, deps.CacheTTL)
			cancel()
		}

		if deps.KafkaPublisher != nil {
			_ = deps.KafkaPublisher.Publish(r.Context(), kafka.Event{
				Kind:   kafka.LinkVisitedKind,
				LinkID: data.linkID,
				UserID: data.ownerID,
				Ts:     time.Now().UTC(),
			})
		}

		http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
	}
}
