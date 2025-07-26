package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sxd0/go_url-shortener/internal/gateway/middleware"
	"github.com/sxd0/go_url-shortener/pkg/req"
	"github.com/sxd0/go_url-shortener/pkg/res"
	"github.com/sxd0/go_url-shortener/proto/gen/go/linkpb"
)

type LinkHandler struct {
	Client linkpb.LinkServiceClient
}

func NewLinkHandler(client linkpb.LinkServiceClient) *LinkHandler {
	return &LinkHandler{Client: client}
}

type CreateLinkRequest struct {
	Url string `json:"url" validate:"required,url"`
}

type UpdateLinkRequest struct {
	Hash   string `json:"hash" validate:"required"`
	NewUrl string `json:"new_url" validate:"required,url"`
}

func (h *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := req.HandleBody[CreateLinkRequest](&w, r)
		if err != nil {
			return
		}

		userID, err := middleware.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		resp, err := h.Client.CreateLink(r.Context(), &linkpb.CreateLinkRequest{
			Url:    payload.Url,
			UserId: uint32(userID),
		})
		if err != nil {
			http.Error(w, "failed to create link: "+err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, resp, http.StatusOK)
	}
}

func (h *LinkHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		resp, err := h.Client.GetAllLinks(r.Context(), &linkpb.GetAllLinksRequest{
			UserId: uint32(userID),
		})
		if err != nil {
			http.Error(w, "failed to get links: "+err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, resp, http.StatusOK)
	}
}

func (h *LinkHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		resp, err := h.Client.GetLinkByHash(r.Context(), &linkpb.GetLinkByHashRequest{
			Hash: hash,
		})
		if err != nil {
			http.Error(w, "failed to get link: "+err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, resp, http.StatusOK)
	}
}

func (h *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := req.HandleBody[UpdateLinkRequest](&w, r)
		if err != nil {
			return
		}

		resp, err := h.Client.UpdateLink(r.Context(), &linkpb.UpdateLinkRequest{
			Hash: payload.Hash,
			Url:  payload.NewUrl,
		})
		if err != nil {
			http.Error(w, "failed to update link: "+err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, resp, http.StatusOK)
	}
}

func (h *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		_, err = h.Client.DeleteLink(r.Context(), &linkpb.DeleteLinkRequest{
			Id: uint32(id),
		})
		if err != nil {
			http.Error(w, "failed to delete link: "+err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, map[string]string{"status": "deleted"}, http.StatusOK)
	}
}

func (h *LinkHandler) DeleteByHash() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		if hash == "" {
			http.Error(w, "invalid hash", http.StatusBadRequest)
			return
		}

		getResp, err := h.Client.GetLinkByHash(r.Context(), &linkpb.GetLinkByHashRequest{
			Hash: hash,
		})
		if err != nil || getResp == nil || getResp.Link == nil {
			http.Error(w, "link not found", http.StatusNotFound)
			return
		}

		_, err = h.Client.DeleteLink(r.Context(), &linkpb.DeleteLinkRequest{
			Id: getResp.Link.Id,
		})
		if err != nil {
			http.Error(w, "failed to delete link: "+err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, map[string]string{"status": "deleted"}, http.StatusOK)
	}
}
