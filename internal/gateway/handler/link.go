package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	httpx "github.com/sxd0/go_url-shortener/internal/gateway/http"
	"github.com/sxd0/go_url-shortener/internal/gateway/middleware"
	"github.com/sxd0/go_url-shortener/proto/gen/go/linkpb"
)

type LinkHandler struct {
	Client linkpb.LinkServiceClient
}

func NewLinkHandler(client linkpb.LinkServiceClient) *LinkHandler {
	return &LinkHandler{Client: client}
}

type createReq struct {
	Url string `json:"url" validate:"required,url"`
}

func (h *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, err := middleware.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		body, err := httpx.Decode[createReq](r)
		if err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := httpx.Validate(body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := h.Client.CreateLink(r.Context(), &linkpb.CreateLinkRequest{
			Url:    body.Url,
			UserId: uint32(uid),
		})
		if err != nil {
			middleware.WriteGRPCError(w, err)
			return
		}
		httpx.JSON(w, resp, http.StatusOK)
	}
}

func (h *LinkHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, err := middleware.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		resp, err := h.Client.GetAllLinks(r.Context(), &linkpb.GetAllLinksRequest{
			UserId: uint32(uid),
			Limit:  100,
			Offset: 0,
		})
		if err != nil {
			middleware.WriteGRPCError(w, err)
			return
		}
		httpx.JSON(w, map[string]any{
			"total": resp.Total,
			"items": resp.Links,
		}, http.StatusOK)
	}
}

func (h *LinkHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		resp, err := h.Client.GetLinkByHash(r.Context(), &linkpb.GetLinkByHashRequest{
			Hash: hash,
		})
		if err != nil {
			middleware.WriteGRPCError(w, err)
			return
		}
		httpx.JSON(w, resp, http.StatusOK)
	}
}

type updateReq struct {
	Id   uint32 `json:"id"`
	Hash string `json:"hash"`
	Url  string `json:"url" validate:"required,url"`
}

func (h *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := httpx.Decode[updateReq](r)
		if err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := httpx.Validate(body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := h.Client.UpdateLink(r.Context(), &linkpb.UpdateLinkRequest{
			Id:   body.Id,
			Url:  body.Url,
			Hash: body.Hash,
		})
		if err != nil {
			middleware.WriteGRPCError(w, err)
			return
		}

		httpx.JSON(w, resp, http.StatusOK)
	}
}

func (h *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		_, err = h.Client.DeleteLink(r.Context(), &linkpb.DeleteLinkRequest{
			Id: uint32(id),
		})
		if err != nil {
			middleware.WriteGRPCError(w, err)
			return
		}
		httpx.JSON(w, map[string]string{"status": "deleted"}, http.StatusOK)
	}
}

func (h *LinkHandler) DeleteByHash() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		getResp, err := h.Client.GetLinkByHash(r.Context(), &linkpb.GetLinkByHashRequest{
			Hash: hash,
		})
		if err != nil || getResp.GetLink() == nil {
			middleware.WriteGRPCError(w, err)
			return
		}

		_, err = h.Client.DeleteLink(r.Context(), &linkpb.DeleteLinkRequest{
			Id: getResp.Link.Id,
		})
		if err != nil {
			middleware.WriteGRPCError(w, err)
			return
		}
		httpx.JSON(w, map[string]string{"status": "deleted"}, http.StatusOK)
	}
}
