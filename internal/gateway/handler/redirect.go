package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sxd0/go_url-shortener/internal/gateway/middleware"
	"github.com/sxd0/go_url-shortener/proto/gen/go/linkpb"
	"github.com/sxd0/go_url-shortener/proto/gen/go/statpb"
)

type RedirectDeps struct {
	LinkClient linkpb.LinkServiceClient
	StatClient statpb.StatServiceClient
}

func RedirectHandler(deps Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		if hash == "" {
			http.NotFound(w, r)
			return
		}

		grpcCtx := middleware.AttachCommonMD(r.Context(), r)

		linkResp, err := deps.LinkClient.GetLinkByHash(grpcCtx, &linkpb.GetLinkByHashRequest{Hash: hash})
		if err != nil || linkResp.GetLink() == nil {
			http.NotFound(w, r)
			return
		}

		ownerID := linkResp.GetLink().GetUserId()
		_, _ = deps.StatClient.AddClick(grpcCtx, &statpb.AddClickRequest{
			LinkId: linkResp.GetLink().GetId(),
			UserId: uint64(ownerID),
		})

		http.Redirect(w, r, linkResp.GetLink().GetUrl(), http.StatusFound)
	}
}
