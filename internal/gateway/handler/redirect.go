package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sxd0/go_url-shortener/proto/gen/go/linkpb"
	"github.com/sxd0/go_url-shortener/proto/gen/go/statpb"
)

type RedirectDeps struct {
	LinkClient linkpb.LinkServiceClient
	StatClient statpb.StatServiceClient
}

func RedirectHandler(deps Deps) http.HandlerFunc {
	linkClient := deps.LinkClient
	statClient := deps.StatClient
	verifier := deps.Verifier

	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		if hash == "" {
			http.NotFound(w, r)
			return
		}

		linkResp, err := linkClient.GetLinkByHash(r.Context(), &linkpb.GetLinkByHashRequest{
			Hash: hash,
		})
		if err != nil || linkResp.GetLink() == nil {
			http.NotFound(w, r)
			return
		}

		var userID uint64 = 0
		if auth := r.Header.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") && verifier != nil {
			token := strings.TrimPrefix(auth, "Bearer ")
			if id, err := verifier.ParseToken(token); err == nil {
				userID = uint64(id)
			}
		} else {
			if v := r.Header.Get("X-User-ID"); v != "" {
				if id, err := strconv.ParseUint(v, 10, 64); err == nil {
					userID = id
				}
			}
		}

		_, _ = statClient.AddClick(r.Context(), &statpb.AddClickRequest{
			LinkId: linkResp.Link.Id,
			UserId: userID,
		})

		http.Redirect(w, r, linkResp.Link.Url, http.StatusFound)
	}
}
