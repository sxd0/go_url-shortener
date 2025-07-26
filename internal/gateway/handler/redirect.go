package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sxd0/go_url-shortener/proto/gen/go/linkpb"
	"github.com/sxd0/go_url-shortener/proto/gen/go/statpb"
)

func RedirectHandler(linkClient linkpb.LinkServiceClient, statClient statpb.StatServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		if hash == "" {
			http.NotFound(w, r)
			return
		}

		linkResp, err := linkClient.GetLinkByHash(r.Context(), &linkpb.GetLinkByHashRequest{Hash: hash})
		if err != nil || linkResp == nil || linkResp.Link == nil {
			http.NotFound(w, r)
			return
		}

		var userID uint64 = 0
		if v := r.Header.Get("X-User-ID"); v != "" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				userID = id
			}
		}

		_, _ = statClient.AddClick(r.Context(), &statpb.AddClickRequest{
			LinkId: linkResp.Link.Id,
			UserId: userID,
		})

		http.Redirect(w, r, linkResp.Link.Url, http.StatusFound)
	}
}
