package handler

import (
	"context"
	"encoding/json"
	"log"
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

type cachedLink struct {
	Dest    string `json:"u"`
	LinkID  uint   `json:"lid"`
	OwnerID uint   `json:"oid"`
}

func RedirectHandler(deps RedirectDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		if hash == "" {
			http.NotFound(w, r)
			return
		}

		cacheKey := "cache:link:v2:" + hash

		if deps.Cache != nil {
			ctx, cancel := context.WithTimeout(r.Context(), 120*time.Millisecond)
			if raw, ok, err := deps.Cache.GetString(ctx, cacheKey); err == nil && ok {
				cancel()
				var cl cachedLink
				if err := json.Unmarshal([]byte(raw), &cl); err == nil && cl.Dest != "" {
					if u, err := url.ParseRequestURI(cl.Dest); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
						if deps.KafkaPublisher != nil {
							_ = deps.KafkaPublisher.Publish(r.Context(), kafka.Event{
								Kind:   kafka.LinkVisitedKind,
								LinkID: cl.LinkID,
								UserID: cl.OwnerID,
								Ts:     time.Now().UTC(),
							})
						}
						http.Redirect(w, r, cl.Dest, http.StatusTemporaryRedirect)
						return
					}
				}
			} else {
				cancel()
			}
		}

		v, err, _ := sf.Do(cacheKey, func() (interface{}, error) {
			grpcCtx := middleware.AttachCommonMD(r.Context(), r)
			resp, err := deps.LinkClient.GetLinkByHash(grpcCtx, &linkpb.GetLinkByHashRequest{Hash: hash})
			if err != nil || resp.GetLink() == nil {
				return nil, err
			}
			lnk := resp.GetLink()
			return cachedLink{
				Dest:    lnk.GetUrl(),
				LinkID:  uint(lnk.GetId()),
				OwnerID: uint(lnk.GetUserId()),
			}, nil
		})
		if err != nil || v == nil {
			http.NotFound(w, r)
			return
		}
		cl := v.(cachedLink)

		u, err := url.ParseRequestURI(cl.Dest)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}

		if deps.Cache != nil {
			if b, err := json.Marshal(cl); err == nil {
				ctx, cancel := context.WithTimeout(r.Context(), 120*time.Millisecond)
				_ = deps.Cache.SetString(ctx, cacheKey, string(b), deps.CacheTTL)
				cancel()
			}
		}

		if deps.KafkaPublisher != nil {
			if err := deps.KafkaPublisher.Publish(r.Context(), kafka.Event{
				Kind:   kafka.LinkVisitedKind,
				LinkID: cl.LinkID,
				UserID: cl.OwnerID,
				Ts:     time.Now().UTC(),
			}); err != nil {
				log.Printf("[KAFKA][PUBLISH][FAIL] link_id=%d user_id=%d err=%v", cl.LinkID, cl.OwnerID, err)
			} else {
				log.Printf("[KAFKA][PUBLISH][OK] link_id=%d user_id=%d", cl.LinkID, cl.OwnerID)
			}
		} else {
			log.Printf("[KAFKA][SKIP] publisher is nil")
		}

		http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
	}
}
